package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

type ReduceClientStockParams struct {
	ClientID       int64       `json:"touser"`
	ProducToReduce []Product   `json:"productoadd"`
	Amount         []int64     `json:"amount"`
	Transaction    Transaction `json:"transaction_id"`
	AfterPaying    func([]map[string]interface{}) error
}

func (store *Store) ReduceClientStockTx(ctx context.Context, arg ReduceClientStockParams) error {
	log.Info().Msg("Reducing Transaction")
	// var result ReduceClientStockResult

	err := store.execTx(ctx, func(q *Queries) error {
		client, err := q.GetUserForUpdate(ctx, arg.ClientID)
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
			// {
			// 	"user_contact": client.PhoneNumber,
			// 	"user_address": client.Address,
			// 	"user_email":   client.Email,
			// },
			// {
			// 	"mpesa_receipt_number": arg.Transaction.MpesaReceiptNumber,
			// 	"phone_number":         arg.Transaction.PhoneNumber,
			// 	"transaction_amount":   arg.Transaction.Amount,
			// 	"status":               "success",
			// },
		}

		var totalStockValue int64
		totalStockValue = 0
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
				"productPrice":    addProduct.UnitPrice,
				"totalBill":       int32(arg.Amount[index]) * addProduct.UnitPrice,
			})

			totalStockValue += int64(int32(arg.Amount[index]) * addProduct.UnitPrice)
		}

		jsonClientProducts, err := json.Marshal(clientProducts)
		if err != nil {
			return err
		}

		_, err = q.UpdateUserStock(ctx, UpdateUserStockParams{
			UserID: client.UserID,
			Stock:  jsonClientProducts,
		})
		if err != nil {
			return err
		}

		// update users stock value minus/reduce
		_, err = q.UpdateUserStockValue(ctx, UpdateUserStockValueParams{
			UserID: int32(client.UserID),
			Value:  -int64(arg.Transaction.Amount),
		})
		if err != nil {
			return err
		}

		return arg.AfterPaying(receiptData)
	})

	return err
}
