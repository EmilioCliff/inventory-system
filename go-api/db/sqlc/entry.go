package db

import (
	"context"
	"time"
)

// type StoreGetEntryByNameRow struct {
// 	IssuedDate         time.Time `json:"issued_date"`
// 	ProductName        string    `json:"product_name"`
// 	TotalProductPrice  int64     `json:"total_product_price"`
// 	TotalQuantityAdded int64     `json:"total_quantity_added"`
// }

type EntryByDate struct {
	IssuedDate   time.Time           `json:"issued_date"`
	Transactions []GetEntryByNameRow `json:"transactions"`
}

func (q *Queries) StoreGetEntryByDate(ctx context.Context) ([]EntryByDate, error) {
	rows, err := q.db.Query(ctx, getEntryByName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store transactions by date
	dateMap := make(map[time.Time][]GetEntryByNameRow)

	for rows.Next() {
		var entry GetEntryByNameRow
		if err := rows.Scan(
			&entry.IssuedDate,
			&entry.ProductName,
			&entry.TotalProductPrice,
			&entry.TotalQuantityAdded,
		); err != nil {
			return nil, err
		}

		if _, ok := dateMap[entry.IssuedDate.Time]; !ok {
			dateMap[entry.IssuedDate.Time] = []GetEntryByNameRow{entry}
		} else {
			dateMap[entry.IssuedDate.Time] = append(dateMap[entry.IssuedDate.Time], entry)
		}
	}

	var result []EntryByDate
	for date, transactions := range dateMap {
		result = append(result, EntryByDate{
			IssuedDate:   date,
			Transactions: transactions,
		})
	}

	return result, nil
}
