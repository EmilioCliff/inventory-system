package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

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

		admin, err := q.GetUserForUpdate(ctx, 1)
		if err != nil {
			return err
		}

		var adminProducts []map[string]interface{}
		if unerr := json.Unmarshal(admin.Stock, &adminProducts); unerr != nil {
			return unerr
		}

		for idx, product := range adminProducts {
			log.Println(product)
			productID, ok := product["productID"].(float64)
			if !ok {
				return fmt.Errorf("Error converting product_id to float64")
			}
			quantity, ok := product["productQuantity"].(float64)
			if !ok {
				return fmt.Errorf("Error converting quantity to float64")
			}

			if int64(productID) == arg.ProductToEdit.ProductID {
				adminProducts[idx] = map[string]interface{}{
					"productID":       result.ProductEdited.ProductID,
					"productName":     result.ProductEdited.ProductName,
					"productQuantity": quantity,
				}
				break
			}
		}

		updatedStock, err := json.Marshal(adminProducts)
		if err != nil {
			return err
		}

		_, err = q.UpdateUserStock(ctx, UpdateUserStockParams{
			UserID: 1,
			Stock:  updatedStock,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
