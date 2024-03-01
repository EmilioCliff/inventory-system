package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	// "google.golang.org/appengine/log"
)

// type SQLStore struct {
// 	*Queries
// 	connPool *pgxpool.Pool
// }

// type Store interface {
// 	Querier
// }

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rberr := tx.Rollback(ctx); rberr != nil {
			return fmt.Errorf("Fn erro: %v rb Error %c", err, rberr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type AddAdminStockParams struct {
	Admin       User    `json:"admin"`
	ProducToAdd Product `json:"productoadd"`
	Amount      int64   `json:"amount"`
}

type AddAdminStockResult struct {
	Admin User
}

func (store *Store) AddAdminStockTx(ctx context.Context, arg AddAdminStockParams) (AddAdminStockResult, error) {
	var result AddAdminStockResult

	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUserForUpdate(ctx, arg.Admin.UserID)
		if err != nil {
			return err
		}

		var userProducts []map[string]interface{}
		if user.Stock != nil {
			if unerr := json.Unmarshal(user.Stock, &userProducts); unerr != nil {
				return unerr
			}
		}

		productExists := false
		for _, product := range userProducts {

			if id, ok := product["productID"].(float64); ok {
				idInt := int64(id)

				if idInt == arg.ProducToAdd.ProductID {
					quantityFloat := product["productQuantity"].(float64)
					quantityInt := int64(quantityFloat)
					product["productQuantity"] = quantityInt + arg.Amount
					productExists = true
					break
				}
			}
		}

		if !productExists {
			newProduct := map[string]interface{}{
				"productID":       arg.ProducToAdd.ProductID,
				"productName":     arg.ProducToAdd.ProductName,
				"productQuantity": arg.Amount,
			}
			userProducts = append(userProducts, newProduct)
		}

		jsonUserProducts, err := json.Marshal(userProducts)
		if err != nil {
			return err
		}
		// arg.Admin.Stock = jsonUserProducts

		result.Admin, err = q.UpdateUserStock(ctx, UpdateUserStockParams{
			UserID: arg.Admin.UserID,
			Stock:  jsonUserProducts,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

type AddClientStockParams struct {
	FromAdmin   User      `json:"fromuser"`
	ToClient    User      `json:"touser"`
	ProducToAdd []Product `json:"productoadd"`
	Amount      []int64   `json:"amount"`
}

type AddClientStockResult struct {
	FromAdmin        User    `json:"fromuser"`
	ToUser           User    `json:"touser"`
	InvoiceGenerated Invoice `json:"invoice"`
}

func (store *Store) AddClientStockTx(ctx context.Context, arg AddClientStockParams) (AddClientStockResult, error) {
	var result AddClientStockResult

	err := store.execTx(ctx, func(q *Queries) error {
		admin, err := q.GetUserForUpdate(ctx, arg.FromAdmin.UserID)
		if err != nil {
			return err
		}

		var adminProducts []map[string]interface{}
		if unerr := json.Unmarshal(admin.Stock, &adminProducts); unerr != nil {
			return unerr
		}

		client, err := q.GetUserForUpdate(ctx, arg.ToClient.UserID)
		if err != nil {
			return err
		}

		var clientProducts []map[string]interface{}
		if client.Stock != nil {
			if unerr := json.Unmarshal(client.Stock, &clientProducts); unerr != nil {
				return unerr
			}
		}

		invoiceData := []map[string]interface{}{
			{
				"user_contact": client.PhoneNumber,
				"user_address": client.Address,
				"user_email":   client.Email,
			},
		}

		for index, addProduct := range arg.ProducToAdd {
			// Reduce Admins Product
			for _, adminProduct := range adminProducts {
				if id, ok := adminProduct["productID"].(float64); ok {
					idInt := int64(id)
					if idInt == addProduct.ProductID {
						quantityFloat := adminProduct["productQuantity"].(float64)
						quantityInt := quantityFloat
						if quantityInt-float64(arg.Amount[index]) < 0 {
							return fmt.Errorf("Not enough in inventory %v - %v to sell %v", adminProduct["productName"], adminProduct["productQuantity"], arg.Amount[index])
						}
						adminProduct["productQuantity"] = quantityInt - float64(arg.Amount[index])
					}
				}
			}

			found := false
			for _, clientProduct := range clientProducts {
				// Add Client's Product
				if id, ok := clientProduct["productID"].(float64); ok {
					idInt := int64(id)
					if idInt == addProduct.ProductID {
						if quantity, ok := clientProduct["productQuantity"].(float64); ok {
							quantityInt := quantity
							clientProduct["productQuantity"] = quantityInt + float64(arg.Amount[index])
							found = true
							break
						}
					}
				}
			}

			// If product not found in client's stock, add it
			if !found {
				clientProducts = append(clientProducts, map[string]interface{}{
					"productID":       addProduct.ProductID,
					"productName":     addProduct.ProductName,
					"productQuantity": arg.Amount[index],
				})
			}

			// Update invoice data
			invoiceData = append(invoiceData, map[string]interface{}{
				"productID":       float64(addProduct.ProductID),
				"productName":     addProduct.ProductName,
				"productQuantity": arg.Amount[index],
				"totalBill":       int32(arg.Amount[index]) * addProduct.UnitPrice,
			})
		}

		jsonAdminProducts, err := json.Marshal(adminProducts)
		if err != nil {
			return err
		}

		result.FromAdmin, err = q.UpdateUserStock(ctx, UpdateUserStockParams{
			UserID: arg.FromAdmin.UserID,
			Stock:  jsonAdminProducts,
		})
		if err != nil {
			return err
		}

		jsonClientProducts, err := json.Marshal(clientProducts)
		if err != nil {
			return err
		}

		result.ToUser, err = q.UpdateUserStock(ctx, UpdateUserStockParams{
			UserID: arg.ToClient.UserID,
			Stock:  jsonClientProducts,
		})
		if err != nil {
			return err
		}

		jsonInvoiceData, err := json.Marshal(invoiceData)
		if err != nil {
			return err
		}
		timestamp := time.Now().Format("20060102150405")
		createdTime := time.Now().Format("2006-01-02")

		invoiceC := map[string]string{
			"invoice_number":   timestamp,
			"created_at":       createdTime,
			"invoice_username": client.Username,
		}

		pdfBytes, err := utils.SetInvoiceVariables(invoiceC, invoiceData)
		if err != nil {
			log.Printf("Error creating invoice %v", err)
			return err
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			result.InvoiceGenerated, err = q.CreateInvoice(ctx, CreateInvoiceParams{
				InvoiceNumber:       timestamp + utils.RandomString(4),
				UserInvoiceID:       int32(client.UserID),
				InvoiceData:         jsonInvoiceData,
				UserInvoiceUsername: client.Username,
				InvoicePdf:          pdfBytes,
			})
		}()

		wg.Wait()

		return err
	})

	go func() {
		config, _ := utils.ReadConfig("../..")
		if err != nil {
			log.Fatal("Could not log config file: ", err)
		}

		emailSender := utils.NewGmailSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)

		emailBody := fmt.Sprintf(`
		<h1>Hello %s</h1>
		<p>We've issued products. Find the invoice attached below</p>
		<h5>Thank You For Choosing Us.</h5>
	`, result.ToUser.Username)

		_ = emailSender.SendMail("Invoice Issued", emailBody, []string{result.ToUser.Email}, nil, nil, "Invoice.pdf", []byte(result.InvoiceGenerated.InvoicePdf))
	}()

	return result, err
}

type EditProductParams struct {
	ProductToEdit Product `json:"producttoedit"`
}

type EditProductResult struct {
	ProductEdited Product `json:"productedited"`
}

func (store *Store) EditStockTx(ctx context.Context, arg EditProductParams) (EditProductResult, error) {
	var result EditProductResult

	err := store.execTx(ctx, func(q *Queries) error {
		product, err := q.GetProductForUpdate(ctx, arg.ProductToEdit.ProductID)
		if err != nil {
			return err
		}

		result.ProductEdited, err = q.UpdateProduct(ctx, UpdateProductParams{
			ProductID:   product.ProductID,
			ProductName: arg.ProductToEdit.ProductName,
			Packsize:    arg.ProductToEdit.Packsize,
			UnitPrice:   arg.ProductToEdit.UnitPrice,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

type EditUserParams struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Username    string `json:"username"`
}

type EditUserResult struct {
	UserEdited User `json:"useredited"`
}

func (store *Store) EditUserTx(ctx context.Context, arg EditUserParams) (EditUserResult, error) {
	var result EditUserResult

	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUserForUpdate(ctx, arg.UserID)
		if err != nil {
			return err
		}

		if arg.Role == "admin" {
			result.UserEdited, err = q.UpdateUserCredentials(ctx, UpdateUserCredentialsParams{
				Password:    user.Password,
				UserID:      user.UserID,
				Email:       arg.Email,
				Username:    arg.Username,
				Address:     arg.Address,
				PhoneNumber: arg.PhoneNumber,
			})
			if err != nil {
				return err
			}
		} else {
			result.UserEdited, err = q.UpdateUserCredentials(ctx, UpdateUserCredentialsParams{
				Password:    arg.Password,
				Email:       user.Email,
				UserID:      user.UserID,
				Username:    user.Username,
				Address:     user.Address,
				PhoneNumber: user.PhoneNumber,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

type ReduceClientStockParams struct {
	Client         User      `json:"touser"`
	ProducToReduce []Product `json:"productoadd"`
	Amount         []int8    `json:"amount"`
}

type ReduceClientStockResult struct {
	Client           User    `json:"touser"`
	ReceiptGenerated Receipt `json:"invoice"`
}

func (store *Store) ReduceClientStockTx(ctx context.Context, arg ReduceClientStockParams) (ReduceClientStockResult, error) {
	var result ReduceClientStockResult

	err := store.execTx(ctx, func(q *Queries) error {
		client, err := q.GetUserForUpdate(ctx, arg.Client.UserID)
		if err != nil {
			return err
		}

		var clientProducts []map[string]interface{}
		if client.Stock != nil {
			if unerr := json.Unmarshal(client.Stock, &clientProducts); unerr != nil {
				return unerr
			}
		}

		// Generate client's receipt
		receiptData := []map[string]interface{}{
			{
				"user_contact": client.PhoneNumber,
				"user_address": client.Address,
				"user_email":   client.Email,
			},
		}

		for index, addProduct := range arg.ProducToReduce {
			// Reduce Clients Product
			for _, clientProduct := range clientProducts {
				if id, ok := clientProduct["productID"].(float64); ok {
					idInt := int64(id)
					if idInt == addProduct.ProductID {
						// Convert product quantity to int8 before subtraction
						quantityFloat := clientProduct["productQuantity"].(float64)
						quantityInt := quantityFloat
						if quantityInt-float64(arg.Amount[index]) < 0 {
							return fmt.Errorf("Not enough in inventory %v - %v to sell %v", clientProduct["productName"], clientProduct["productQuantity"], arg.Amount[index])
						}
						clientProduct["productQuantity"] = quantityInt - float64(arg.Amount[index])
					}
				}
			}

			// Update receipt data
			receiptData = append(receiptData, map[string]interface{}{
				"productID":       float64(addProduct.ProductID),
				"productName":     addProduct.ProductName,
				"productQuantity": arg.Amount[index],
				"totalBill":       int32(arg.Amount[index]) * addProduct.UnitPrice,
			})
		}

		jsonClientProducts, err := json.Marshal(clientProducts)
		if err != nil {
			return err
		}

		result.Client, err = q.UpdateUserStock(ctx, UpdateUserStockParams{
			UserID: arg.Client.UserID,
			Stock:  jsonClientProducts,
		})
		if err != nil {
			return err
		}

		jsonReceiptData, err := json.Marshal(receiptData)
		if err != nil {
			return err
		}

		timestamp := time.Now().Format("20060102150405")
		createdTime := time.Now().Format("2006-01-02")
		receiptC := map[string]string{
			"receipt_number":   timestamp,
			"created_at":       createdTime,
			"receipt_username": client.Username,
		}

		pdfBytes, err := utils.SetReceiptVariables(receiptC, receiptData)
		if err != nil {
			log.Printf("Error creating receipt pdf %v", err)
			return err
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			result.ReceiptGenerated, err = q.CreateReceipt(ctx, CreateReceiptParams{
				ReceiptNumber:       timestamp + utils.RandomString(5),
				UserReceiptID:       int32(client.UserID),
				ReceiptData:         jsonReceiptData,
				UserReceiptUsername: client.Username,
				ReceiptPdf:          pdfBytes,
			})
		}()

		wg.Wait()

		return err
	})

	go func() {
		config, _ := utils.ReadConfig("../..")
		if err != nil {
			log.Fatal("Could not log config file: ", err)
		}

		emailSender := utils.NewGmailSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)

		emailBody := fmt.Sprintf(`
		<h1>Hello %s</h1>
		<p>We've received your payment. Find the receipt attached below</p>
		<h5>Thank You For Choosing Us.</h5>
	`, result.Client.Username)

		_ = emailSender.SendMail("Receipt Issued", emailBody, []string{result.Client.Email}, nil, nil, "Receipt.pdf", []byte(result.ReceiptGenerated.ReceiptPdf))
	}()
	return result, err
}
