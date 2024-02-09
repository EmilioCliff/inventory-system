package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
}

func randomInt(max, min int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	k := len(letters)
	var word strings.Builder

	for i := 0; i < n; i++ {
		c := letters[rand.Intn(k)]
		word.WriteByte(c)
	}

	return word.String()
}

func RandomName() string {
	return RandomString(rand.Intn(12-5+1) + 5)
}

func RandomPassword() string {
	return RandomString(rand.Intn(10-8+1) + 1)
}

func RandomEmail() string {
	return RandomString(5) + "@gmail.com"
}

func RandomPhoneNumber() string {
	return strconv.Itoa(randomInt(999999999, 100000000))
}

func RandomRole() string {
	roles := []string{"admin", "client"}
	k := len(roles)
	return roles[rand.Intn(k)]
}

func RandomInvoiceReceiptNumber() string {
	return strconv.Itoa(rand.Intn(1000))
}
