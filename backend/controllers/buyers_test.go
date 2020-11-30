package controllers_test

import (
	"backend/models"
	"backend/test_utils"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadBuyers(t *testing.T) {
	nBuyers, _ := test_utils.RandomSliceOfBuyers(5)
	for _, nBuyer := range nBuyers {
		buyersMap[nBuyer.Id] = nBuyer
	}

	jsonBody, err := json.Marshal(nBuyers)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "products.csv")
	if err != nil {
		t.Fail()
	}
	part.Write(jsonBody)

	writer.WriteField("date", date.Format("2006-01-02"))
	writer.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/buyers", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

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
		assert.Equal(t, buyersMap[id].Date.Format("2006-01-02"), storedBuyer.Date.Format("2006-01-02"))
		storedBuyersMap[id].Date = buyersMap[id].Date
	}

	assert.Equal(t, buyersMap, storedBuyersMap)
}

func TestGetBuyers(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/buyers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var retrievedBuyers []*models.Buyer
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedBuyers); err != nil {
		t.Fail()
	}

	retrievedBuyersMap := test_utils.MapOfBuyersFromSlice(retrievedBuyers)

	for id, storedBuyer := range retrievedBuyersMap {
		assert.Equal(t, buyersMap[id].Date.Format("2006-01-02"), storedBuyer.Date.Format("2006-01-02"))
		retrievedBuyersMap[id].Date = buyersMap[id].Date
	}

	assert.Equal(t, retrievedBuyersMap, buyersMap)

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
		transactions[id] = &models.Transaction{
			Id:       id,
			Buyer:    buyer[0],
			Ip:       ip,
			Device:   device,
			Products: transactionProducts,
			DType:    []string{"Transaction"},
		}
	}
}

func TestGetBuyer(t *testing.T) {
	insertTransactionsWithBuyers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", ("/buyers/" + buyerT.Id), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp models.BuyerDetailResponse

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fail()
	}

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

func TestGetBuyerNotFound(t *testing.T) {
	insertTransactionsWithBuyers(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", ("/buyers/" + buyerT.Id + "abc"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
