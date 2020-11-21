package models_test

import (
	"backend/models"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func randomProduct() *models.Product {
	return models.NewProduct(randomString(8), randomString(40), randomInt(3000))
}


func TestInsertManyProducts(t *testing.T){
	products := make([]*models.Product, 5)
	for i:=0; i < 5; i++{
		products[i] = randomProduct()
	}
	models.InsertManyProducts(products)


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

	c, err := models.NewClient()
	if err != nil{
		t.Fail()
	}
	txn := c.NewTxn()

	resp, err := txn.Query(context.Background(), query)
	if err != nil{
		t.Fail()
	}

	var storedProducts struct {
		All []*models.Product
	}

	if err := json.Unmarshal(resp.GetJson(), &storedProducts); err != nil{
		t.Fail()
	}

	assert.Equal(t, products, storedProducts.All)

}