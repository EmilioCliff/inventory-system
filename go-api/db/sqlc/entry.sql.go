// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: entry.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (
    product_name, product_price, quantity_added
) VALUES (
    $1, $2, $3
)
RETURNING entry_id, product_name, product_price, quantity_added, created_at
`

type CreateEntryParams struct {
	ProductName   string `json:"product_name"`
	ProductPrice  int32  `json:"product_price"`
	QuantityAdded int32  `json:"quantity_added"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRow(ctx, createEntry, arg.ProductName, arg.ProductPrice, arg.QuantityAdded)
	var i Entry
	err := row.Scan(
		&i.EntryID,
		&i.ProductName,
		&i.ProductPrice,
		&i.QuantityAdded,
		&i.CreatedAt,
	)
	return i, err
}

const getEntryByName = `-- name: GetEntryByName :many
SELECT
    DATE_TRUNC('day', created_at)::timestamp AS issued_date,
    product_name,
    SUM(product_price) AS total_product_price,
    SUM(quantity_added) AS total_quantity_added
FROM
    entries
GROUP BY
    issued_date, product_name
ORDER BY
    issued_date, product_name
`

type GetEntryByNameRow struct {
	IssuedDate         pgtype.Timestamp `json:"issued_date"`
	ProductName        string           `json:"product_name"`
	TotalProductPrice  int64            `json:"total_product_price"`
	TotalQuantityAdded int64            `json:"total_quantity_added"`
}

func (q *Queries) GetEntryByName(ctx context.Context) ([]GetEntryByNameRow, error) {
	rows, err := q.db.Query(ctx, getEntryByName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEntryByNameRow{}
	for rows.Next() {
		var i GetEntryByNameRow
		if err := rows.Scan(
			&i.IssuedDate,
			&i.ProductName,
			&i.TotalProductPrice,
			&i.TotalQuantityAdded,
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