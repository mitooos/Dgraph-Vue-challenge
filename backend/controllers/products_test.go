package controllers_test

import (
	"backend/models"
	"backend/test_utils"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)


func productsToCsv(products []*models.Product)string{
	ans := ""
	for _, product := range products {
		ans += product.Id + "'" + product.Name + "'" + strconv.Itoa(product.Price) + "\n"
	}
	return ans
}

func TestUploadProducts(t *testing.T){
	nProducts, _ := test_utils.RandomSliceOfProducts(5)
	for _, nProduct := range nProducts{
		productsMap[nProduct.Id] = nProduct
	}
	
	csvBody := []byte(productsToCsv(nProducts))


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

	resp, err := models.ExecuteQuery(query)
	if err != nil{
		t.Fail()
	}

	var storedProducts struct {
		All []*models.Product
	}

	if err := json.Unmarshal(resp.GetJson(), &storedProducts); err != nil{
		t.Fail()
	}

	storedProductsMap := test_utils.MapOfProductsFromSlice(storedProducts.All)

	assert.Equal(t, productsMap, storedProductsMap)
}