package api

import (
	"net/http"
	"time"

	"github.com/EmilioCliff/inventory-system/reports"
	"github.com/gin-gonic/gin"
)

type downloadUserReportsRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

func (server *Server) downloadUserReports(ctx *gin.Context) {
	var req downloadUserReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	toDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	excelBytes, err := server.reportMaker.GenerateUserExcel(ctx, reports.ReportsPayload{
		FromDate: fromDate,
		ToDate:   toDate,
	})

	ctx.JSON(http.StatusOK, gin.H{"data": excelBytes})

	// ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// ctx.Header("Content-Disposition", "attachment; filename=report.xlsx")
	// ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}

func (server *Server) downloadAdminReports(ctx *gin.Context) {
	var req downloadUserReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	toDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	excelBytes, err := server.reportMaker.GenerateManagerReports(ctx, reports.ReportsPayload{
		FromDate: fromDate,
		ToDate:   toDate,
	})

	ctx.JSON(http.StatusOK, gin.H{"data": excelBytes})

	// ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// ctx.Header("Content-Disposition", "attachment; filename=report.xlsx")
	// ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBytes)
}
