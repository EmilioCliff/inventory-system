package db

import (
	"context"
	"encoding/json"
	"fmt"
)

type AddClientStockParams struct {
	FromAdmin    User      `json:"fromuser"`
	ToClient     User      `json:"touser"`
	ProducToAdd  []Product `json:"productoadd"`
	Amount       []int64   `json:"amount"`
	AfterProcess func([]map[string]interface{}) error
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

		var totalStockValue int64
		totalStockValue = 0
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

			totalStockValue += int64(int32(arg.Amount[index]) * addProduct.UnitPrice)
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

		// update users stock value
		q.UpdateUserStockValue(ctx, UpdateUserStockValueParams{
			UserID: int32(client.UserID),
			Value:  totalStockValue,
		})

		return arg.AfterProcess(invoiceData)
	})

	return result, err
}
