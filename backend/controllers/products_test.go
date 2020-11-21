package controllers_test

import (
	"backend/models"
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)


func randomProduct() *models.Product {
	return models.NewProduct(randomString(8), randomString(40), randomInt(3000))
}

func productsToCsv(products []*models.Product)string{
	ans := ""
	for _, product := range products {
		ans += product.Id + "'" + product.Name + "'" + strconv.Itoa(product.Price) + "\n"
	}
	return ans
}

func TestUploadProducts(t *testing.T){
	products := make([]*models.Product, 5)
	productsMap := make(map[string] *models.Product)
	for i:=0; i < 5; i++{
		randomProduct := randomProduct()
		products[i] = randomProduct
		productsMap[randomProduct.Id] = randomProduct
	}

	csvBody := []byte(productsToCsv(products))

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)


	part, err := writer.CreateFormFile("file", "products.csv")
	if err != nil {
		t.Fail()
	}
	part.Write(csvBody)
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/products", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)


	const query = `
	{
		All(func: type(Product)){
			id
			name
			price
			dgraph.type
		}
	}
	`

	c, err := models.NewClient()
	if err != nil{
		t.Fail()
	}
	txn := c.NewTxn()

	resp, err := txn.Query(context.Background(), query)
	if err != nil{
		t.Fail()
	}

	var storedProducts struct {
		All []*models.Product
	}

	if err := json.Unmarshal(resp.GetJson(), &storedProducts); err != nil{
		t.Fail()
	}

	storedProductsMap := make(map[string]*models.Product)
	for _, storedProduct := range storedProducts.All{
		storedProductsMap[storedProduct.Id] = storedProduct
	}

	assert.Equal(t, productsMap, storedProductsMap)
}