package models_test

import (
	"backend/models"
	"backend/test_utils"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestInsertManyProducts(t *testing.T){
	nProducts, _ := test_utils.RandomSliceOfProducts(5)
	models.InsertManyProducts(nProducts)

	for _, nProduct := range nProducts{
		productsMap[nProduct.Id] = nProduct
	}


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