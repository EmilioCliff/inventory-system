package reports

import (
	"context"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
)

type ReportStore struct {
	dbStore *db.Store
}

type ReportMaker interface {
	GetUserInvoiceSummary(ctx context.Context, payload ReportSummaryData) ([]GetUserInvoiceSummaryResponse, error)
	GetUserReceiptSummary(ctx context.Context, payload ReportSummaryData) ([]GetUserReceiptSummaryResponse, error)
	GetAdminPurchaseHistory(ctx context.Context, payload ReportSummaryData) ([]GetAdminPurchaseHistoryResponse, error)
	GetUserHistorySummary(ctx context.Context, payload ReportSummaryData) (GetAdminSalesHistoryResponse, error)

	GenerateUserExcel(ctx context.Context, payload ReportsPayload) ([]byte, error)
	GenerateManagerReports(ctx context.Context, payload ReportsPayload) ([]byte, error)
}

func NewReportMaker(store *db.Store) ReportMaker {
	return &ReportStore{
		dbStore: store,
	}
}
