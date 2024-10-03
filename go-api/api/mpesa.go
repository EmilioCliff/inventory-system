package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/gin-gonic/gin"
)

const (
	registerURL = "https://api.safaricom.co.ke/mpesa/c2b/v1/registerurl"
)

type registerUrlRequest struct {
	ShortCode       string `json:"ShortCode"`
	ResponseType    string `json:"ResponseType"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ValidationURL   string `json:"validationURL"`
}

func (s *Server) registerUrl(ctx *gin.Context) {
	var req registerUrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	accesstoken, err := utils.GenerateAccessToken(s.config.CONSUMER_KEY, s.config.CONSUMER_SECRET)
	if err != nil {
		log.Println("Could not generate access token")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	req.ShortCode = s.config.MPESA_SHORT_CODE
	req.ResponseType = "Completed"

	reqBytes, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	request, err := http.NewRequest(http.MethodPost, registerURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	request.Header.Set("Authorization", "Bearer "+accesstoken)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	var registerResponseBody map[string]interface{}

	if err := json.Unmarshal(responseBody, &registerResponseBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))

		return
	}

	log.Println("register url hit: ", registerResponseBody)

	rspCode, ok := registerResponseBody["ResponseCode"].(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "not expected response code"})

		return
	}

	if rspCode != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "response code failed", "code": rspCode})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": registerResponseBody})
}

type completeTransactionRequest struct {
	TransactionType   string `json:"TransactionType"`
	TransID           string `json:"TransID"`
	TransTime         string `json:"TransTime"`
	TransAmount       string `json:"TransAmount"`
	BusinessShortCode string `json:"BusinessShortCode"`
	BillRefNumber     string `json:"BillRefNumber"`
	InvoiceNumber     string `json:"InvoiceNumber"`
	OrgAccountBalance string `json:"OrgAccountBalance"`
	ThirdPartyTransID string `json:"ThirdPartyTransID"`
	MSISDN            string `json:"MSISDN"`
	FirstName         string `json:"FirstName"`
	MiddleName        string `json:"MiddleName"`
	LastName          string `json:"LastName"`
}

func (s *Server) completeTransaction(ctx *gin.Context) {
	var req completeTransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("failed to bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fullname := fmt.Sprintf("%s %s %s", req.FirstName, req.MiddleName, req.LastName)

	amount, err := strconv.Atoi(req.TransAmount)
	if err != nil {
		log.Println("failed to convert transAmount to int: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	orgAmount, err := strconv.Atoi(req.OrgAccountBalance)
	if err != nil {
		log.Println("failed to convert org account balance to int: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// YYYYMMDDHHmmss

	transactionTime, err := time.Parse("20060102150405", req.TransTime)
	if err != nil {
		log.Println("failed to parse transaction time: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := s.store.CreateC2BTransaction(ctx, db.CreateC2BTransactionParams{
		Fullname:          fullname,
		Phone:             req.MSISDN,
		Amount:            int64(amount),
		TransactionID:     req.TransID,
		OrgAccountBalance: int64(orgAmount),
		TransactionTime:   transactionTime,
	})
	if err != nil {
		log.Println("failed to create transaction: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("complete url hit: ", req)
	log.Println("transaction: ", transaction)

	ctx.JSON(http.StatusOK, gin.H{
		"ResultCode": 0,
		"ResultDesc": "Accepted",
	})
}

func (s *Server) listC2BTransactions(ctx *gin.Context) {
	transactions, err := s.store.ListC2BTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// {
// 	"BillRefNumber": "iiiii",
// 	"BusinessShortCode": "600426",
// 	"FirstName": "John",
// 	"InvoiceNumber": "",
// 	"LastName": "",
// 	"MSISDN": "254708374149",
// 	"MiddleName": "Doe",
// 	"OrgAccountBalance": "5490845.42",
// 	"ThirdPartyTransID": "",
// 	"TransAmount": "1000.00",
// 	"TransID": "SJ242OOZWW",
// 	"TransTime": "20241002202238",
// 	"TransactionType": "Pay Bill"
//   }

func (s *Server) validateTransaction(ctx *gin.Context) {
	var req any

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("validate url hit: ", req)

	ctx.JSON(http.StatusOK, gin.H{"ResultCode": "0", "ResultDesc": "Accepted"})
}
