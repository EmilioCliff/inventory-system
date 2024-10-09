package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/gin-gonic/gin"
)

const (
	registerURL = "https://api.safaricom.co.ke/mpesa/c2b/v1/registerurl"
)

// 9090757

type registerUrlRequest struct {
	ShortCode       string `json:"ShortCode"`
	ResponseType    string `json:"ResponseType"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ValidationURL   string `json:"ValidationURL"`
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

func (s *Server) completeTransaction(ctx *gin.Context) {
	var rq any

	if err := ctx.ShouldBindJSON(&rq); err != nil {
		log.Println("failed to bind json: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req, _ := rq.(map[string]interface{})

	fullname := fmt.Sprintf("%s", req["FirstName"])

	// YYYYMMDDHHmmss

	transactionTime, err := time.Parse("20060102150405", req["TransTime"].(string))
	if err != nil {
		log.Println("failed to parse transaction time: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := s.store.CreateC2BTransaction(ctx, db.CreateC2BTransactionParams{
		Fullname:          fullname,
		Phone:             "****",
		Amount:            req["TransAmount"].(string),
		TransactionID:     req["TransID"].(string),
		OrgAccountBalance: req["OrgAccountBalance"].(string),
		TransactionTime:   transactionTime,
	})
	if err != nil {
		log.Println("failed to create transaction: ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("transaction: ", transaction)

	ctx.JSON(http.StatusOK, gin.H{
		"ResultCode": 0,
		"ResultDesc": "Accepted",
	})
}

func (s *Server) validateTransaction(ctx *gin.Context) {
	var req any

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Println("validate url hit: ", req)

	ctx.JSON(http.StatusOK, gin.H{"ResultCode": "0", "ResultDesc": "Accepted"})
}
