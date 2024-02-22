package db

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomUserTest() (CreateUserParams, User, error) {
	password, err := utils.GeneratePasswordHash(utils.RandomString(6))
	arg := CreateUserParams{
		Username:    utils.RandomName(),
		Password:    password,
		Email:       utils.RandomEmail(),
		PhoneNumber: utils.RandomPhoneNumber(),
		Address:     utils.RandomName(),
		Role:        utils.RandomRole(),
		Stock:       nil,
	}

	// _, product, _ := createProductTest()

	// Initialize stock data with a random product
	// stockData := []map[string]interface{}{
	// 	{
	// 		"productID":       product.ProductID,
	// 		"productName":     product.ProductName,
	// 		"productQuantity": 10,
	// 	},
	// }
	// stockData := []map[string]interface{}{
	// 	{
	// 		"productID":       9,
	// 		"productName":     "HIV testkit 100",
	// 		"productQuantity": 10,
	// 	},
	// }

	// jsonStockData, err := json.Marshal(stockData)
	// if err != nil {
	// 	fmt.Println("Error marshaling invoice data:", err)
	// }
	// arg.Stock = jsonStockData
	user, err := testStore.CreateUser(context.Background(), arg)

	return arg, user, err
}

func TestCreateUser(t *testing.T) {
	arg, user, err := CreateRandomUserTest()
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)
	require.Equal(t, arg.Address, user.Address)
	require.Equal(t, arg.Role, user.Role)

	require.NotZero(t, user.UserID)
	require.NotZero(t, user.CreatedAt)
}

func TestDeleteUser(t *testing.T) {
	_, user, err := CreateRandomUserTest()
	// user, err := testStore.GetUser(context.Background(), 4)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	err = testStore.DeleteUser(context.Background(), user.UserID)
	require.NoError(t, err)

	_, err = testStore.GetUser(context.Background(), user.UserID)
	require.Error(t, err)
}

func TestGetUser(t *testing.T) {
	_, user, err := CreateRandomUserTest()

	require.NoError(t, err)
	require.NotEmpty(t, user)

	var newUser User

	newUser, err = testStore.GetUser(context.Background(), user.UserID)

	require.Equal(t, user.Username, newUser.Username)
	require.Equal(t, user.Password, newUser.Password)
	require.Equal(t, user.Email, newUser.Email)
	require.Equal(t, user.PhoneNumber, newUser.PhoneNumber)
	require.Equal(t, user.Address, newUser.Address)
	require.Equal(t, user.Stock, newUser.Stock)
	require.Equal(t, user.Role, newUser.Role)
}

func TestUpdateUserStock(t *testing.T) {
	_, newUser, err := CreateRandomUserTest()

	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	user, err := testStore.GetUser(context.Background(), newUser.UserID)
	if err != nil {
		fmt.Println(err)
	}

	var prevStock []map[string]interface{}
	if prevStock != nil {
		err = json.Unmarshal(user.Stock, &prevStock)
		if err != nil {
			fmt.Println("Error unmarshaling: ", err)
		}
	}

	for _, product := range prevStock {
		if id, ok := product["productID"].(float64); ok && id == 4 {
			if quantity, ok := product["productQuantity"].(float64); ok {
				product["productQuantity"] = quantity + 2
			}
		}
	}
	prevStockJSON, err := json.Marshal(prevStock)
	if err != nil {
		fmt.Println("Error marshaling invoice data:", err)
	}
	updateArg := UpdateUserStockParams{
		UserID: user.UserID,
		Stock:  prevStockJSON,
	}
	updatedUser, updatedErr := testStore.UpdateUserStock(context.Background(), updateArg)
	if updatedErr != nil {
		fmt.Println(err)
	}
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Password, updatedUser.Password)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.PhoneNumber, updatedUser.PhoneNumber)
	require.Equal(t, user.Address, updatedUser.Address)
	// require.NotEqual(t, user.Stock, updatedUser.Stock)
	require.Equal(t, user.Role, updatedUser.Role)
}

func TestListUser(t *testing.T) {
	users, err := testStore.ListUser(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)
	for _, user := range users {
		require.NotEmpty(t, user.UserID)
		require.NotEmpty(t, user.Username)
		require.NotEmpty(t, user.Email)
		require.NotEmpty(t, user.PhoneNumber)
		require.NotEmpty(t, user.Address)
		require.NotEmpty(t, user.Role)
		require.NotEmpty(t, user.Password)
		// require.NotEmpty(t, user.Stock)
		// fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v", user.UserID, user.Username, user.Email, user.PhoneNumber, user.Address, user.Role)
	}
}
