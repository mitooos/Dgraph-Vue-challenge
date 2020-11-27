package controllers_test

import (
	"backend/models"
	"backend/test_utils"
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomIp()string{
	return fmt.Sprintf("%d.%d.%d.%d", test_utils.RandomInt(256), test_utils.RandomInt(256), test_utils.RandomInt(256), test_utils.RandomInt(256))
}

func randomMapOfTransactions(size int)(map[string]*models.Transaction, error){
	ans := make(map[string]*models.Transaction)
	for i := 0; i < size; i++{
		buyer, _ := test_utils.RandomSliceOfBuyers(1)
		transactionProducts, _ := test_utils.RandomSliceOfProducts(5)
		id := test_utils.RandomString(12)
		ip := randomIp()
		device := test_utils.RandomString(10)
		
		if err := models.InsertManyBuyers(buyer); err != nil{
			return nil, err
		}
		buyers = append(buyers, buyer...)
		buyersMap[buyer[0].Id] = buyer[0]

		if err := models.InsertManyProducts(transactionProducts); err != nil{
			return nil, err
		}
		products = append(products, transactionProducts...)
		for _, product := range transactionProducts{
			productsMap[product.Id] = product
		}

		ans[id] = &models.Transaction{
			Id: id,
			Buyer: buyer[0],
			Ip: ip,
			Device: device,
			Products: transactionProducts,
			DType: []string{"Transaction"},
		}
	}
	return ans, nil
}

func transactionsToFileFormat(transactions map[string]*models.Transaction)string{
	ans := ""
	for _, transaction := range transactions {
		productIds := ""
		for _, product := range transaction.Products {
			productIds += (product.Id + ",")
		}
		//remove last comma
		productIds = productIds[0:len(productIds) - 1]
		ans += fmt.Sprintf("#%s\000%s\000%s\000%s\000(%s)\000\000", transaction.Id, transaction.Buyer.Id, transaction.Ip, transaction.Device, productIds)
	}
	return ans
}

func TestInsertTransactions(t *testing.T){
	transactionsFileFormat := transactionsToFileFormat(transactions)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)


	part, err := writer.CreateFormFile("file", "transactions")
	if err != nil {
		t.Fail()
	}
	part.Write([]byte(transactionsFileFormat))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/transactions", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	query := `
	{
		All(func: type(Transaction)){
			id
			date
			ip
			device
			dgraph.type
			buyer{
				id
				name
				age
				dgraph.type
			}
			products{
				id
				name
				price
				dgraph.type
			}
		}
	}
	`

	resp, err := models.ExecuteQuery(query)
	if err != nil{
		t.Fail()
	}

	var storedTransactions struct {
		All []*models.Transaction
	}

	if err := json.Unmarshal(resp.GetJson(), &storedTransactions); err != nil{
		t.Fail()
	}

	for _, storedTransaction := range storedTransactions.All{
		transaction := transactions[storedTransaction.Id]
		transaction.Date = storedTransaction.Date
		assert.Equal(t, test_utils.MapOfProductsFromSlice(transaction.Products), test_utils.MapOfProductsFromSlice(storedTransaction.Products))
		storedTransaction.Products = transaction.Products
		assert.Equal(t, transaction, storedTransaction)
	}
}