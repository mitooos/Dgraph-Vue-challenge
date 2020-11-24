package models_test

import (
	"backend/models"
	"backend/test_utils"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestInsertManyBuyers(t *testing.T){
	models.InsertManyBuyers(buyers)

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