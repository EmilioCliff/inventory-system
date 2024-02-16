package db

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/stretchr/testify/require"
)

// var store = NewStore(testConnPool)

func TestAddAdminStockTx(t *testing.T) {
	store := NewStore(testConnPool)

	_, product, _ := createProductTest()
	arg := CreateUserParams{
		Username:    utils.RandomName(),
		Password:    utils.RandomPassword(),
		Email:       utils.RandomEmail(),
		PhoneNumber: utils.RandomPhoneNumber(),
		Address:     utils.RandomName(),
		Role:        utils.RandomRole(),
		Stock:       nil,
	}

	// initialAmount := int64(10)

	// stockData := []map[string]interface{}{
	// 	{
	// 		"productID":       product.ProductID,
	// 		"productName":     product.ProductName,
	// 		"productQuantity": 10,
	// 	},
	// }

	// jsonStockData, err := json.Marshal(stockData)
	// if err != nil {
	// 	fmt.Println("Error marshaling invoice data:", err)
	// }
	// arg.Stock = jsonStockData
	adminAccount, err := testStore.CreateUser(context.Background(), arg)

	// Unmarshal the stock data
	var initialAdminStock []map[string]interface{}
	if adminAccount.Stock != nil {
		err = json.Unmarshal(adminAccount.Stock, &initialAdminStock)
		if err != nil {
			fmt.Println("Error unmarshaling data")
		}
	}

	amountToAdd := int64(10)

	n := 2
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	resCh := make(chan AddAdminStockResult, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			updatedAdmin, err := store.AddAdminStockTx(context.Background(), AddAdminStockParams{
				Admin:       adminAccount,
				ProducToAdd: product,
				Amount:      amountToAdd,
			})
			errCh <- err
			resCh <- updatedAdmin
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(resCh)
	}()

	for i := 0; i < n; i++ {

		updatedAdmin := <-resCh
		err := <-errCh

		require.NotEmpty(t, updatedAdmin)
		require.NoError(t, err)
		require.Equal(t, adminAccount.UserID, updatedAdmin.Admin.UserID)
		require.Equal(t, adminAccount.PhoneNumber, updatedAdmin.Admin.PhoneNumber)

		var updatedAdminStock []map[string]interface{}
		err = json.Unmarshal(updatedAdmin.Admin.Stock, &updatedAdminStock)

		for _, updatedproduct := range updatedAdminStock {
			if id, ok := updatedproduct["productID"].(float64); ok {
				idInt := int64(id)
				if idInt == product.ProductID {
					updatedAmountFloat := updatedproduct["productQuantity"].(float64)
					updatedAmountInt := int64(updatedAmountFloat)
					// require.NotEqual(t, initialAmount, updatedAmountInt)
					require.Equal(t, updatedAmountInt, amountToAdd*int64(i+1))
					// require.Equal(t, updatedAmountInt-initialAmount*amountToAdd, amountToAdd)
				}
			}
		}
	}
}

func TestAddClientStockTx(t *testing.T) {
	store := NewStore(testConnPool)

	_, adminUser, _ := CreateRandomUserTest()
	_, clientUser, _ := CreateRandomUserTest()

	productsToAdd := []Product{
		{
			ProductID:   1,
			ProductName: "Test Product 1",
			UnitPrice:   100,
		},
		{
			ProductID:   2,
			ProductName: "Test Product 2",
			UnitPrice:   200,
		},
	}

	amounts := []int64{2, 4}

	adminStockData := []map[string]interface{}{
		{
			"productID":       1,
			"productName":     "Test Product 1",
			"productQuantity": 200,
		},
		{
			"productID":       2,
			"productName":     "Test Product 2",
			"productQuantity": 300,
		},
	}

	jsonAdminStockData, _ := json.Marshal(adminStockData)

	_, _ = store.UpdateUserStock(context.Background(), UpdateUserStockParams{
		UserID: adminUser.UserID,
		Stock:  jsonAdminStockData,
	})

	addClientStockParams := AddClientStockParams{
		FromAdmin:   adminUser,
		ToClient:    clientUser,
		ProducToAdd: productsToAdd,
		Amount:      amounts,
	}
	n := 5
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	resCh := make(chan AddClientStockResult, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := store.AddClientStockTx(context.Background(), addClientStockParams)

			errCh <- err
			resCh <- result
		}()
	}

	go func() {
		wg.Wait()

		close(errCh)
		close(resCh)
	}()

	for i := 0; i < n; i++ {
		fmt.Println(i)
		err := <-errCh
		result := <-resCh

		require.NoError(t, err)
		require.NotEmpty(t, result)

		adminStockAfterUpdate, _ := store.GetUser(context.Background(), adminUser.UserID)
		var adminStockAfterUpdateData []map[string]interface{}
		_ = json.Unmarshal(adminStockAfterUpdate.Stock, &adminStockAfterUpdateData)

		clientStockAfterUpdate, _ := store.GetUser(context.Background(), clientUser.UserID)
		var clientStockAfterUpdateData []map[string]interface{}
		_ = json.Unmarshal(clientStockAfterUpdate.Stock, &clientStockAfterUpdateData)

		for index, product := range productsToAdd {
			for _, stock := range adminStockAfterUpdateData {
				if id, ok := stock["productID"].(float64); ok && int64(id) == product.ProductID {
					initialQuantity := adminStockData[index]["productQuantity"].(int)
					expectedAdminStockQuantity := initialQuantity - (int(amounts[index]) * (i))
					require.NotEmpty(t, expectedAdminStockQuantity)
					// require.Equal(t, expectedAdminStockQuantity, int(stock["productQuantity"].(float64)))
				}
			}

			found := false
			for _, stock := range clientStockAfterUpdateData {
				if id, ok := stock["productID"].(float64); ok && int64(id) == product.ProductID {
					// require.Equal(t, (i+1)*int(amounts[index]), int(stock["productQuantity"].(float64)))
					found = true
					break
				}
			}
			require.True(t, found, "Product not found in client's stock")
		}

		require.NotNil(t, result.InvoiceGenerated)
		require.NotEmpty(t, result.InvoiceGenerated.InvoiceNumber)
		require.Equal(t, clientUser.UserID, int64(result.InvoiceGenerated.UserInvoiceID))
	}
}

func TestEditStockTx(t *testing.T) {
	store := NewStore(testConnPool)

	product, err := store.CreateProduct(context.Background(), CreateProductParams{
		ProductName: "Test Product",
		Packsize:    "50 Kits",
		UnitPrice:   200,
	})
	require.NoError(t, err)
	require.NotEmpty(t, product)

	product.ProductName = "Edited Product"
	product.Packsize = "100 Kits"
	product.UnitPrice = 400
	// Edit the product
	editedProduct, err := store.EditStockTx(context.Background(), EditProductParams{
		ProductToEdit: product,
	})
	require.NoError(t, err)
	require.NotEmpty(t, product)

	// Verify the edited product
	require.NotEmpty(t, editedProduct.ProductEdited)
	require.Equal(t, product.ProductID, editedProduct.ProductEdited.ProductID)
	require.Equal(t, product.ProductName, editedProduct.ProductEdited.ProductName)
	require.Equal(t, product.Packsize, editedProduct.ProductEdited.Packsize)
	require.Equal(t, product.UnitPrice, editedProduct.ProductEdited.UnitPrice)
}

func TestReduceClientStockTx(t *testing.T) {
	store := NewStore(testConnPool)

	_, clientUser, _ := CreateRandomUserTest()

	productsToReduce := []Product{
		{
			ProductID:   1,
			ProductName: "Test Product 1",
			UnitPrice:   100,
		},
		{
			ProductID:   2,
			ProductName: "Test Product 2",
			UnitPrice:   200,
		},
	}

	amounts := []int8{5, 5}

	clientStockData := []map[string]interface{}{
		{
			"productID":       1,
			"productName":     "Test Product 1",
			"productQuantity": 100,
		},
		{
			"productID":       2,
			"productName":     "Test Product 2",
			"productQuantity": 100,
		},
	}

	jsonClientStockData, _ := json.Marshal(clientStockData)

	_, _ = store.UpdateUserStock(context.Background(), UpdateUserStockParams{
		UserID: clientUser.UserID,
		Stock:  jsonClientStockData,
	})

	reduceClientStockParams := ReduceClientStockParams{
		Client:         clientUser,
		ProducToReduce: productsToReduce,
		Amount:         amounts,
	}

	n := 5
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	resCh := make(chan ReduceClientStockResult, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := store.ReduceClientStockTx(context.Background(), reduceClientStockParams)

			errCh <- err
			resCh <- result
		}()
	}

	go func() {
		wg.Wait()

		close(errCh)
		close(resCh)
	}()

	for i := 0; i < n; i++ {
		err := <-errCh
		result := <-resCh

		require.NoError(t, err)
		require.NotEmpty(t, result)

		clientStockAfterReduction, _ := store.GetUser(context.Background(), clientUser.UserID)
		var clientStockAfterReductionData []map[string]interface{}
		_ = json.Unmarshal(clientStockAfterReduction.Stock, &clientStockAfterReductionData)

		for _, product := range productsToReduce {
			for _, stock := range clientStockAfterReductionData {
				if id, ok := stock["productID"].(float64); ok && int64(id) == product.ProductID {
					// require.Equal(t, clientStockData[index]["productQuantity"].(int)-(int(amounts[index])*(i+1)), int(stock["productQuantity"].(float64)))
				}
			}
		}

		require.NotNil(t, result.ReceiptGenerated)
		require.NotEmpty(t, result.ReceiptGenerated.ReceiptNumber)
		require.Equal(t, clientUser.UserID, int64(result.ReceiptGenerated.UserReceiptID))
	}
}
