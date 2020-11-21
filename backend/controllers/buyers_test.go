package controllers_test

import (
	"backend/models"
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func randomBuyer() *models.Buyer {
	return models.NewBuyer(randomString(8), randomString(20), randomInt(100))
}

func TestUploadBuyers(t *testing.T){
	buyers := make([]*models.Buyer, 5)
	buyersMap := make(map[string]*models.Buyer)
	for i:=0; i < 5; i++{
		randomBuyer := randomBuyer()
		buyers[i] = randomBuyer
		buyersMap[randomBuyer.Id] = randomBuyer
	}

	jsonBody, err := json.Marshal(buyers)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)


	part, err := writer.CreateFormFile("file", "products.csv")
	if err != nil {
		t.Fail()
	}
	part.Write(jsonBody)
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

	var storedBuyers struct {
		All []*models.Buyer
	}

	if err := json.Unmarshal(resp.GetJson(), &storedBuyers); err != nil{
		t.Fail()
	}

	storedBuyersMap := make(map[string]*models.Buyer)
	for _, storedBuyer := range storedBuyers.All{
		storedBuyersMap[storedBuyer.Id] = storedBuyer
	}

	assert.Equal(t, buyersMap, storedBuyersMap)
}