package models_test

import (
	"backend/models"
	"backend/test_utils"
	"os"
	"testing"
)

var buyers []*models.Buyer
var buyersMap map[string] *models.Buyer
var products []*models.Product
var productsMap map[string]*models.Product

func TestMain(m *testing.M){
	models.LoadSchemas()

	buyers, buyersMap = test_utils.RandomSliceOfBuyers(5)
	products, productsMap = test_utils.RandomSliceOfProducts(5)

	code := m.Run()

	os.Exit(code)
}