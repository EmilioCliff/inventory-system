package db

import (
	"context"
	"time"
)

type StoreGetInvoicesByDateRow struct {
	IssuedDate  time.Time `json:"issued_date"`
	NumInvoices int64     `json:"num_invoices"`
	InvoiceData []byte    `json:"invoice_data"`
}

func (q *Queries) StoreGetInvoicesByDate(ctx context.Context) ([]StoreGetInvoicesByDateRow, error) {
	rows, err := q.db.Query(ctx, getInvoicesByDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []StoreGetInvoicesByDateRow{}
	for rows.Next() {
		var i StoreGetInvoicesByDateRow
		if err := rows.Scan(&i.IssuedDate, &i.NumInvoices, &i.InvoiceData); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type StoreGetReceiptsByDateRow struct {
	IssuedDate  time.Time `json:"issued_date"`
	NumReceipts int64     `json:"num_receipts"`
	ReceiptData []byte    `json:"receipt_data"`
}

func (q *Queries) StoreGetReceiptsByDate(ctx context.Context) ([]StoreGetReceiptsByDateRow, error) {
	rows, err := q.db.Query(ctx, getReceiptsByDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []StoreGetReceiptsByDateRow{}
	for rows.Next() {
		var i StoreGetReceiptsByDateRow
		if err := rows.Scan(&i.IssuedDate, &i.NumReceipts, &i.ReceiptData); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type StoreGetUserInvoicesByDateRow struct {
	IssuedDate  time.Time `json:"issued_date"`
	NumInvoices int64     `json:"num_invoices"`
	InvoiceData []byte    `json:"invoice_data"`
}

func (q *Queries) StoreGetUserInvoicesByDate(ctx context.Context, userInvoiceID int32) ([]StoreGetUserInvoicesByDateRow, error) {
	rows, err := q.db.Query(ctx, getUserInvoicesByDate, userInvoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []StoreGetUserInvoicesByDateRow{}
	for rows.Next() {
		var i StoreGetUserInvoicesByDateRow
		if err := rows.Scan(&i.IssuedDate, &i.NumInvoices, &i.InvoiceData); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type StoreGetUserReceiptsByDateRow struct {
	IssuedDate  time.Time `json:"issued_date"`
	NumReceipts int64     `json:"num_receipts"`
	ReceiptData []byte    `json:"receipt_data"`
}

func (q *Queries) StoreGetUserReceiptsByDate(ctx context.Context, userReceiptID int32) ([]StoreGetUserReceiptsByDateRow, error) {
	rows, err := q.db.Query(ctx, getUserReceiptsByDate, userReceiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []StoreGetUserReceiptsByDateRow{}
	for rows.Next() {
		var i StoreGetUserReceiptsByDateRow
		if err := rows.Scan(&i.IssuedDate, &i.NumReceipts, &i.ReceiptData); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
