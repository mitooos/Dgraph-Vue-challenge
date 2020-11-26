package models_test

import (
	"backend/models"
	"backend/test_utils"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestInsertManyBuyers(t *testing.T){
	nBuyers, _ := test_utils.RandomSliceOfBuyers(5)
	models.InsertManyBuyers(nBuyers)

	for _, nBuyer := range nBuyers{
		buyersMap[nBuyer.Id] = nBuyer
	}

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
	retrievedBuyers, err := models.GetBuyers()

	retrievedBuyerssMap := test_utils.MapOfBuyersFromSlice(retrievedBuyers)

	assert.NoError(t, err)
	assert.Equal(t, buyersMap, retrievedBuyerssMap)
}