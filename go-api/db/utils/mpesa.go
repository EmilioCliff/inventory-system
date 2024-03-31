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
	sandbox      = "174379"
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

	accessToken, err := generateAccessToken(consumerKey, consumerSecret)
	if err != nil {
		fmt.Println("Error generating access token:", err)
		return transactionID, err
	}

	// newNumber := setPhoneNumber(phoneNumber)

	// callback := fmt.Sprintf("https://4676-105-163-2-242.ngrok-free.app/transaction/%v%v", transactionID, fmt.Sprintf("%03d", userID))
	callback := fmt.Sprintf("https://hip-letters-production.up.railway.app/transaction/%v%v", transactionID, fmt.Sprintf("%03d", userID))
	requestBody := map[string]interface{}{
		"BusinessShortCode": sandbox,
		"Password":          generatePassword(sandbox, config.PASSKEY),
		"Timestamp":         time.Now().Format("20060102150405"),
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            amount,
		"PartyA":            "254718750145",
		"PartyB":            sandbox,
		"PhoneNumber":       "254718750145",
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

	log.Println(stkResponseBody)

	return transactionID, nil
}

func setPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) == 10 && phoneNumber[0] == '0' {
		phoneNumber = "254" + phoneNumber[1:]
	}
	return phoneNumber
}

// var accessToken string
// var expiryTimestamp time.Time

// func getAccessToken(consumerKey, consumerSecret string) (string, error) {
// 	// if time.Now().Before(expiryTimestamp) {
// 	// 	return accessToken, nil
// 	// }

// 	newToken, err := generateAccessToken(consumerKey, consumerSecret)
// 	if err != nil {
// 		return "", err
// 	}

// 	// accessToken = newToken
// 	// expiryTimestamp = time.Now().Add(3600 * time.Second)

// 	return accessToken, nil
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

	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", err
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("Access token not found in response")
	}
	log.Println(accessToken)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Request failed with status code %d", resp.StatusCode)
	}

	return accessToken, nil
}

func generatePassword(shortCode string, key string) string {
	password := shortCode + key + time.Now().Format("20060102150405")
	return base64.StdEncoding.EncodeToString([]byte(password))
}
