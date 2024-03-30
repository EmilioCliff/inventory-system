package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/EmilioCliff/inventory-system/db/utils"
)

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
		config, err := utils.ReadConfig("../..")
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
