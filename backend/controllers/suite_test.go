package controllers_test

import (
	"backend/controllers"
	"backend/models"
	"backend/test_utils"
	"os"
	"testing"

	"github.com/go-chi/chi"
)

var router *chi.Mux

var buyers []*models.Buyer
var buyersMap map[string] *models.Buyer
var products []*models.Product
var productsMap map[string]*models.Product
var transactions map [string]*models.Transaction


func TestMain(m *testing.M){
	models.LoadSchemas()
	router = controllers.Router()

	populateDb()

	code := m.Run()

	os.Exit(code)
}

func populateDb(){
	buyers, buyersMap = test_utils.RandomSliceOfBuyers(5)
	products, productsMap = test_utils.RandomSliceOfProducts(5)
	
	var err error
	transactions, err = randomMapOfTransactions(10)
	if err != nil {
		panic(err)
	}

	models.InsertManyBuyers(buyers)
	models.InsertManyProducts(products)
}