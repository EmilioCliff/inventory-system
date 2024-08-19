package db

import (
	"context"
	"encoding/json"
)

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

		if arg.Amount != 0 {
			_, err = q.CreateEntry(ctx, CreateEntryParams{
				ProductID:     int32(arg.ProducToAdd.ProductID),
				ProductName:   arg.ProducToAdd.ProductName,
				ProductPrice:  arg.ProducToAdd.UnitPrice * int32(arg.Amount),
				QuantityAdded: int32(arg.Amount),
			})
			if err != nil {
				return err
			}
		}

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
