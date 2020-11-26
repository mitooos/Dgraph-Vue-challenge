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



func TestUploadBuyers(t *testing.T){
	nBuyers, _ := test_utils.RandomSliceOfBuyers(5)
	for _, nBuyer := range nBuyers{
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

	resp, err := models.ExecuteQuery(query)
	if err != nil{
		t.Fail()
	}

	var storedBuyers struct {
		All []*models.Buyer
	}

	if err := json.Unmarshal(resp.GetJson(), &storedBuyers); err != nil{
		t.Fail()
	}

	storedBuyersMap := test_utils.MapOfBuyersFromSlice(storedBuyers.All)

	assert.Equal(t, buyersMap, storedBuyersMap)
}

func TestGetBuyers(t *testing.T){
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/buyers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var retrievedBuyers []*models.Buyer
	if err := json.Unmarshal(w.Body.Bytes(), &retrievedBuyers); err != nil{
		t.Fail()
	}

	retrievedBuyersMap := test_utils.MapOfBuyersFromSlice(retrievedBuyers)

	assert.Equal(t, retrievedBuyersMap, buyersMap)

}