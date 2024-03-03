package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL      = "https://sandbox.safaricom.co.ke"
	shortCode    = "174379"
	lipaEndpoint = "/mpesa/stkpush/v1/processrequest"
	callbackPath = "/callback"
)

var transactionID string

func SendSTK(amount string, userID int64, phoneNumber string) (string, error) {
	config, err := ReadConfig("../..")
	if err != nil {
		log.Fatal("Could not log config file: ", err)
	}

	consumerKey := config.CONSUMER_KEY
	consumerSecret := config.CONSUMER_SECRET

	transactionID = time.Now().Format("20060102150405")

	log.Println("Generate Token")
	accessToken, err := generateAccessToken(consumerKey, consumerSecret)
	if err != nil {
		return transactionID, err
	}

	// callback := fmt.Sprintf("https://e864-105-163-157-51.ngrok-free.app/transaction/%v%v", transactionID, fmt.Sprintf("%03d", userID))
	callback := fmt.Sprintf("https://hip-letters-production.up.railway.app/transaction/%v%v", transactionID, fmt.Sprintf("%03d", userID))
	log.Println(callback)
	requestBody := map[string]interface{}{
		"BusinessShortCode": shortCode,
		"Password":          generatePassword(shortCode, config.PASSKEY),
		"Timestamp":         time.Now().Format("20060102150405"),
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            amount,
		"PartyA":            phoneNumber,
		"PartyB":            shortCode,
		"PhoneNumber":       phoneNumber,
		"CallBackURL":       callback,
		"AccountReference":  "Cliff Test",
		"TransactionDesc":   "Pay Bob For Test",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return transactionID, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL+lipaEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return transactionID, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return transactionID, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return transactionID, err
	}

	var stkResponseBody map[string]interface{}
	err = json.Unmarshal(responseBody, &stkResponseBody)
	if err != nil {
		return transactionID, err
	}
	// log.Println(stkResponseBody)

	// if stkResponseBody["ResponseCode"] != 0 {
	// 	return transactionID, fmt.Errorf("Unsuccessful STK push")
	// }

	// log.Println("STK Push Response Body:", string(responseBody))
	return transactionID, nil
}

// func mpesaCallbackHandler(c *gin.Context) {

// 	body, err := io.ReadAll(c.Request.Body)
// 	if err != nil {
// 		fmt.Println("Error reading request body:", err)
// 		c.String(http.StatusInternalServerError, "Internal Server Error")
// 		return
// 	}

// 	var callbackBody map[string]interface{}
// 	err = json.Unmarshal(body, &callbackBody)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling callback data:", err)
// 		return
// 	}

// 	log.Println(callbackBody)
// 	fmt.Println(callbackBody)

// 	merchantRequestID, _ := callbackBody["MerchantRequestID"].(string)
// 	checkoutRequestID, _ := callbackBody["CheckoutRequestID"].(string)
// 	resultCode, _ := callbackBody["ResultCode"].(float64)
// 	resultDesc, _ := callbackBody["ResultDesc"].(string)

// 	fmt.Println("MerchantRequestID:", merchantRequestID)
// 	fmt.Println("CheckoutRequestID:", checkoutRequestID)
// 	fmt.Println("ResultCode:", resultCode)
// 	fmt.Println("ResultDesc:", resultDesc)

// 	c.JSON(http.StatusOK, gin.H{"Successful": "Reached"})
// }

func generateAccessToken(consumerKey string, consumerSecret string) (string, error) {
	authString := consumerKey + ":" + consumerSecret
	encodedAuthString := base64.StdEncoding.EncodeToString([]byte(authString))

	url := baseURL + "/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+encodedAuthString)

	log.Println("Sending request")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Print the response body
	log.Println("Response Body:", string(body))

	var tokenResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("Access token not found in response")
	}
	log.Println(accessToken)

	return accessToken, nil
}

func generatePassword(shortCode string, key string) string {
	passkey := key
	timestamp := time.Now().Format("20060102150405")

	password := shortCode + passkey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(password))
}
