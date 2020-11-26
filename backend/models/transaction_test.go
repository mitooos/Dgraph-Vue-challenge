package models_test

import (
	"backend/models"
	"backend/test_utils"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestInsertTransaction(t *testing.T){
	buyer, _ := test_utils.RandomSliceOfBuyers(1)
	transactionProducts, transactionProductsMap := test_utils.RandomSliceOfProducts(5)
	id := test_utils.RandomString(8)
	date := time.Now()
	buyerId := buyer[0].Id
	ip := test_utils.RandomString(16)
	device := test_utils.RandomString(10)

	productsIds := ""
	for _, product := range transactionProducts{
		productsIds += (product.Id + " ")
		products = append(products, product)
	}
	productsIds = productsIds[0:len(productsIds) - 1]

	buyers = append(buyers, buyer[0])
	if err := models.InsertManyBuyers(buyer); err != nil{
		t.Fail()
	}
	if err := models.InsertManyProducts(transactionProducts); err != nil{
		t.Fail()
	}

	if  err := models.InsertTransaction(id, date, buyerId, ip, device, productsIds); err != nil{
		panic(err)
		t.Fail()
	}
	
	transaction := &models.Transaction{
		Id: id,
		Date: date,
		Buyer: buyer[0],
		Ip: ip,
		Device: device,
		Products: transactionProducts,
		DType: []string{"Transaction"},
	}


	query := fmt.Sprintf(`
	{
		Transaction(func: eq(id, "%s" )){
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
	`, id)

	resp, err := models.ExecuteQuery(query)
	if err != nil{
		t.Fail()
	}

	var storedTransaction struct {
		Transaction []*models.Transaction
	}

	if err := json.Unmarshal(resp.Json, &storedTransaction); err != nil{
		t.Fail()
	}

	storedProductsMap := test_utils.MapOfProductsFromSlice(storedTransaction.Transaction[0].Products)

	assert.Equal(t, transactionProductsMap, storedProductsMap)
	assert.Equal(t, transaction.Date.Format("2006-01-02T15:04:05"), storedTransaction.Transaction[0].Date.Format("2006-01-02T15:04:05"))

	storedTransaction.Transaction[0].Date = transaction.Date
	storedTransaction.Transaction[0].Products = transaction.Products
	assert.Equal(t, transaction, storedTransaction.Transaction[0])
}