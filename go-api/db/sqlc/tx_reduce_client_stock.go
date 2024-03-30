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

type ReduceClientStockParams struct {
	Client         User        `json:"touser"`
	ProducToReduce []Product   `json:"productoadd"`
	Amount         []int8      `json:"amount"`
	Transaction    Transaction `json:"transaction_id"`
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

		receiptData := []map[string]interface{}{
			{
				"user_contact": client.PhoneNumber,
				"user_address": client.Address,
				"user_email":   client.Email,
			},
		}

		for index, addProduct := range arg.ProducToReduce {
			for _, clientProduct := range clientProducts {
				if id, ok := clientProduct["productID"].(float64); ok {
					idInt := int64(id)
					if idInt == addProduct.ProductID {
						quantityFloat := clientProduct["productQuantity"].(float64)
						quantityInt := quantityFloat
						if quantityInt-float64(arg.Amount[index]) < 0 {
							return fmt.Errorf("Not enough in inventory %v - %v to sell %v", clientProduct["productName"], clientProduct["productQuantity"], arg.Amount[index])
						}
						clientProduct["productQuantity"] = quantityInt - float64(arg.Amount[index])
					}
				}
			}

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

		timestamp := arg.Transaction.TransactionID
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
				ReceiptNumber:       timestamp,
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
