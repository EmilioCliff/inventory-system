package db

import (
	"context"
	"encoding/json"
	"fmt"
)

type ReduceClientStockParams struct {
	Client         User        `json:"touser"`
	ProducToReduce []Product   `json:"productoadd"`
	Amount         []int8      `json:"amount"`
	Transaction    Transaction `json:"transaction_id"`
	AfterPaying    func() error
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

		return arg.AfterPaying()
	})

	return result, err
}
