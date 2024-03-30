package db

import "context"

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
