// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
)

type Querier interface {
	CreateInvoice(ctx context.Context, arg CreateInvoiceParams) (Invoice, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateReceipt(ctx context.Context, arg CreateReceiptParams) (Receipt, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteProduct(ctx context.Context, productID int64) error
	DeleteUser(ctx context.Context, userID int64) error
	GetInvoice(ctx context.Context, invoiceNumber string) (Invoice, error)
	GetInvoiceByID(ctx context.Context, invoiceID int64) (Invoice, error)
	GetProduct(ctx context.Context, productID int64) (Product, error)
	GetProductForUpdate(ctx context.Context, productID int64) (Product, error)
	GetReceipt(ctx context.Context, receiptNumber string) (Receipt, error)
	GetReceiptByID(ctx context.Context, receiptID int64) (Receipt, error)
	GetUser(ctx context.Context, userID int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByUsename(ctx context.Context, username string) (User, error)
	GetUserForUpdate(ctx context.Context, userID int64) (User, error)
	GetUserInvoicesByID(ctx context.Context, userInvoiceID int32) ([]Invoice, error)
	GetUserInvoicesByUsername(ctx context.Context, userInvoiceUsername string) ([]Invoice, error)
	GetUserReceiptsByID(ctx context.Context, userReceiptID int32) ([]Receipt, error)
	GetUserReceiptsByUsername(ctx context.Context, userReceiptUsername string) ([]Receipt, error)
	ListInvoices(ctx context.Context) ([]Invoice, error)
	ListProduct(ctx context.Context) ([]Product, error)
	ListReceipts(ctx context.Context) ([]Receipt, error)
	ListUser(ctx context.Context) ([]User, error)
	SearchILikeProducts(ctx context.Context, productName string) ([]Product, error)
	SearchILikeUsers(ctx context.Context, username string) ([]string, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateUserCredentials(ctx context.Context, arg UpdateUserCredentialsParams) (User, error)
	UpdateUserPasswordFisrtLogin(ctx context.Context, arg UpdateUserPasswordFisrtLoginParams) (User, error)
	UpdateUserStock(ctx context.Context, arg UpdateUserStockParams) (User, error)
}

var _ Querier = (*Queries)(nil)