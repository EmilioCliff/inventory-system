package api

// cache route lists
const (
	UserProducts               = "/users/products/"
	ListProducts               = "/products/"
	ListAllProducts            = "/allproducts/"
	GetUser                    = "/users/"
	ListUsers                  = "/users/admin"
	UserInvoices               = "/users/invoices/"
	UserReceipts               = "/users/receipts/"
	ListInvoices               = "/invoices/admin"
	GetInvoice                 = "/invoices/"
	ListReceipts               = "/receipts/admin"
	GetReceipt                 = "/receipts/"
	AllTransactions            = "/transactions/all"
	SuccesfulTransactions      = "/transactions/successfull"
	FailedTransactions         = "/transactions/failed"
	UserTransactions           = "/user/transactions/all/"
	UserSuccessfulTransactions = "/user/transactions/successful/"
	UserFailedTransactions     = "/user/transactions/failed/"
	GetUserTransaction         = "/user/transactions/"
	UserReceivedHistory        = "/history/received/"
	AllUserReceiverHistory     = "/history/all_received"
	UserSoldHistory            = "/history/sold/"
	UserDebt                   = "/history/debt/"
	AllUserDebt                = "/history/all_debt"
	AdminHistory               = "/history/admin"
	TestGroup                  = "/history/test"
)