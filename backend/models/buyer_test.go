package models_test

import (
	"backend/models"
	"backend/test_utils"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertManyBuyers(t *testing.T) {
	nBuyers, _ := test_utils.RandomSliceOfBuyers(5)
	models.InsertManyBuyers(nBuyers)

	for _, nBuyer := range nBuyers {
		buyersMap[nBuyer.Id] = nBuyer
	}

	const query = `
	{
		All(func: type(Buyer)){
			id
			name
			age
			date
			dgraph.type
		}
	}
	`

	resp, err := models.ExecuteQuery(query)
	if err != nil {
		t.Fail()
	}

	var storedBuyers struct {
		All []*models.Buyer
	}

	if err := json.Unmarshal(resp.GetJson(), &storedBuyers); err != nil {
		t.Fail()
	}

	storedBuyersMap := test_utils.MapOfBuyersFromSlice(storedBuyers.All)

	for id, storedBuyer := range storedBuyersMap {
		assert.Equal(t, buyersMap[id].Date.Format("2006-01-02T15:04:05"), storedBuyer.Date.Format("2006-01-02T15:04:05"))
		storedBuyersMap[id].Date = buyersMap[id].Date
	}

	assert.Equal(t, buyersMap, storedBuyersMap)
}

func TestGetBuyers(t *testing.T) {
	retrievedBuyers, err := models.GetBuyers("1000", "0")

	retrievedBuyersMap := test_utils.MapOfBuyersFromSlice(retrievedBuyers)

	for id, storedBuyer := range retrievedBuyersMap {
		assert.Equal(t, buyersMap[id].Date.Format("2006-01-02T15:04:05"), storedBuyer.Date.Format("2006-01-02T15:04:05"))
		retrievedBuyersMap[id].Date = buyersMap[id].Date
	}

	assert.NoError(t, err)
	assert.Equal(t, buyersMap, retrievedBuyersMap)
}

var buyersWithSameIp map[string]*models.Buyer
var buyerT *models.Buyer
var mostBoughtProducts map[string]*models.Product

func insertTransactionsWithBuyers(t *testing.T) {
	ip := test_utils.RandomString(16)
	buyersWithSameIp = make(map[string]*models.Buyer)
	mostBoughtProducts = make(map[string]*models.Product)

	transactionProducts, _ := test_utils.RandomSliceOfProducts(10)

	productsIds := ""
	for _, product := range transactionProducts {
		productsIds += (product.Id + " ")
		products = append(products, product)
		mostBoughtProducts[product.Id] = product
	}
	productsIds = productsIds[0 : len(productsIds)-1]
	if err := models.InsertManyProducts(transactionProducts); err != nil {
		t.Fail()
	}

	for i := 0; i < 5; i++ {
		buyer, _ := test_utils.RandomSliceOfBuyers(1)
		id := test_utils.RandomString(8)
		date := time.Now()
		buyerId := buyer[0].Id
		device := test_utils.RandomString(10)

		buyers = append(buyers, buyer[0])
		buyersWithSameIp[buyer[0].Id] = buyer[0]
		buyersMap[buyer[0].Id] = buyer[0]
		if err := models.InsertManyBuyers(buyer); err != nil {
			t.Fail()
		}

		for _, transactionProduct := range transactionProducts {
			productsMap[transactionProduct.Id] = transactionProduct
		}

		if err := models.InsertTransaction(id, date, buyerId, ip, device, productsIds); err != nil {
			panic(err)
			t.Fail()
		}
		buyerT = buyer[0]
	}
}

func TestGetBuyer(t *testing.T) {
	insertTransactionsWithBuyers(t)
	resp, err := models.GetBuyer(buyerT.Id)
	assert.Nil(t, err)

	assert.Equal(t, buyerT.Date.Format("2006-01-02"), resp.Buyer[0].Date.Format("2006-01-02"))
	resp.Buyer[0].Date = buyerT.Date

	resp.Buyer[0].Transactions = []models.Transaction(nil)
	assert.Equal(t, buyerT, resp.Buyer[0])

	buyersWithSameIpMap := test_utils.MapOfBuyersFromSlice(resp.BuyersWithSameIP)

	for id, buyerWithSameIP := range buyersWithSameIpMap {
		assert.Equal(t, buyersWithSameIp[id].Date.Format("2006-01-02:04"), buyerWithSameIP.Date.Format("2006-01-02:04"))
		buyersWithSameIpMap[id].Date = buyersWithSameIp[id].Date
	}

	delete(buyersWithSameIp, buyerT.Id)
	assert.Equal(t, buyersWithSameIp, buyersWithSameIpMap)

	for _, mostBoughtProduct := range resp.RecommendedProducts {
		assert.Equal(t, 5, mostBoughtProduct.Count)
		assert.NotNil(t, mostBoughtProducts[mostBoughtProduct.Id])
	}
}
