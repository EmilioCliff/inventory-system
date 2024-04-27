package db

import (
	"context"
	"time"
)

type StoreGetEntryByNameRow struct {
	IssuedDate    time.Time `json:"issued_date"`
	NumEntries    int64     `json:"num_entries"`
	ProductName   string    `json:"product_name"`
	ProductPrice  int32     `json:"product_price"`
	QuantityAdded int32     `json:"quantity_added"`
}

func (q *Queries) StoreGetEntryByName(ctx context.Context) ([]StoreGetEntryByNameRow, error) {
	rows, err := q.db.Query(ctx, getEntryByName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []StoreGetEntryByNameRow{}
	for rows.Next() {
		var i StoreGetEntryByNameRow
		if err := rows.Scan(
			&i.IssuedDate,
			&i.NumEntries,
			&i.ProductName,
			&i.ProductPrice,
			&i.QuantityAdded,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
