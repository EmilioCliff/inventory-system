package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	shortCode = "7169782"
)

var transactionID string

func SendSTK(amount string, userID int64, phoneNumber string) (string, string, error) {
	config, err := ReadConfig("../..")
	if err != nil {
		log.Fatal().Msgf("Could not log config file: %v", err)
	}

	consumerKey := config.CONSUMER_KEY
	consumerSecret := config.CONSUMER_SECRET

	transactionID = time.Now().Format("20060102150405")

	accessToken, err := generateAccessToken(consumerKey, consumerSecret)
	if err != nil {
		fmt.Println("Error generating access token:", err)
		return "", transactionID, err
	}

	newNumber := setPhoneNumber(phoneNumber)
	log.Info().Msgf("number: %v", newNumber)

	callback := fmt.Sprintf("https://secretive-window-production.up.railway.app/transaction/%v%v", transactionID, fmt.Sprintf("%03d", userID))
	requestBody := map[string]interface{}{
		"BusinessShortCode": shortCode,
		"Password":          generatePassword(shortCode, config.PASSKEY),
		"Timestamp":         time.Now().Format("20060102150405"),
		"TransactionType":   "CustomerBuyGoodsOnline",
		"Amount":            amount,
		"PartyA":            newNumber,
		"PartyB":            "9090757",
		"PhoneNumber":       newNumber,
		"CallBackURL":       callback,
		"AccountReference":  "Kokomed Supplies",
		"TransactionDesc":   "Pay Sold Products",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", transactionID, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.safaricom.co.ke/mpesa/stkpush/v1/processrequest", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", transactionID, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", transactionID, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", transactionID, err
	}

	var stkResponseBody map[string]interface{}
	err = json.Unmarshal(responseBody, &stkResponseBody)
	if err != nil {
		return "", transactionID, err
	}

	log.Info().Msgf("stkResponseBody: %v", stkResponseBody)

	description, ok := stkResponseBody["ResponseDescription"].(string)
	if !ok {
		return "failed to parse mpesa metaData", transactionID, nil
	}

	return description, transactionID, nil
}

func setPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) == 10 && phoneNumber[0] == '0' {
		phoneNumber = "254" + phoneNumber[1:]
	}
	return phoneNumber
}

func generateAccessToken(consumerKey string, consumerSecret string) (string, error) {
	authString := consumerKey + " : " + consumerSecret
	encodedAuthString := base64.StdEncoding.EncodeToString([]byte(authString))
	log.Debug().Msgf("Encoded Auth: %v", encodedAuthString)

	url := "https://api.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error Number 0:%s", err))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+encodedAuthString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error Number 1: %s", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error Number 2: %s", err))
	}

	var tokenResponse map[string]interface{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error Number 3: %s", err))
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("Access token not found in response")
	}
	log.Info().Msgf("tokenResponse: %v", tokenResponse)

	return accessToken, nil
}

func generatePassword(shortCode string, key string) string {
	password := shortCode + key + time.Now().Format("20060102150405")
	return base64.StdEncoding.EncodeToString([]byte(password))
}
