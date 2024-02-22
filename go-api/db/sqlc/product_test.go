package db

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/stretchr/testify/require"
)

func createProductTest() (CreateProductParams, Product, error) {
	// num := utils.RandomInvoiceReceiptNumber()
	price, err := strconv.ParseInt(utils.RandomInvoiceReceiptNumber(), 10, 32)
	if err != nil {
		fmt.Println("Erro converting str to int32: ", err)
	}

	arg := CreateProductParams{
		// ProductName: fmt.Sprintf("HIV testkit %v", num),
		ProductName: "HIV testkit 100",
		UnitPrice:   int32(price),
		Packsize:    "52 testkits",
	}
	product, err := testStore.CreateProduct(context.Background(), arg)
	return arg, product, err
}
func TestCreateProduct(t *testing.T) {
	arg, product, err := createProductTest()
	require.NoError(t, err)
	require.NotEmpty(t, product)
	require.NotEmpty(t, arg)

	require.Equal(t, arg.ProductName, product.ProductName)
	require.Equal(t, arg.UnitPrice, product.UnitPrice)
	require.Equal(t, arg.Packsize, product.Packsize)
	require.NotZero(t, product.ProductID)
	require.NotZero(t, product.CreatedAt)
}

func TestDeleteProduct(t *testing.T) {
	// _, product, err := createProductTest()
	// require.NoError(t, err)
	// require.NotEmpty(t, product)

	// err = testStore.DeleteProduct(context.Background(), product.ProductID)
	err := testStore.DeleteProduct(context.Background(), 8)
	require.NoError(t, err)
}

func TestGetProduct(t *testing.T) {
	_, product, err := createProductTest()
	require.NoError(t, err)
	require.NotEmpty(t, product)

	productGet, err := testStore.GetProduct(context.Background(), product.ProductID)
	require.NoError(t, err)
	require.NotEmpty(t, productGet)

	require.Equal(t, product.ProductName, productGet.ProductName)
	require.Equal(t, product.UnitPrice, productGet.UnitPrice)
	require.Equal(t, product.Packsize, productGet.Packsize)
	require.Equal(t, product.ProductID, productGet.ProductID)
}

func TestListProduct(t *testing.T) {
	products, err := testStore.ListProduct(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, products)

	for _, product := range products {
		require.NotEmpty(t, product.ProductID)
		require.NotEmpty(t, product.ProductName)
		require.NotEmpty(t, product.UnitPrice)
		require.NotEmpty(t, product.Packsize)
		require.NotEmpty(t, product.CreatedAt)
		// fmt.Printf("%v\t%v\t%v\t%v\t%v", product.ProductID, product.ProductName, product.UnitPrice, product.Packsize, product.CreatedAt)
	}
}

func TestUpdateProductPrice(t *testing.T) {
	_, product, err := createProductTest()
	require.NoError(t, err)
	require.NotEmpty(t, product)

	arg := UpdateProductParams{
		ProductID:   product.ProductID,
		UnitPrice:   10,
		ProductName: "New",
		Packsize:    "100",
	}

	updatedProduct, err := testStore.UpdateProduct(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedProduct)

	require.NotEqual(t, product.Packsize, updatedProduct.Packsize)
	require.Equal(t, product.ProductID, updatedProduct.ProductID)
	require.NotEqual(t, product.ProductName, updatedProduct.ProductName)
	require.NotEqual(t, product.UnitPrice, updatedProduct.UnitPrice)
}
