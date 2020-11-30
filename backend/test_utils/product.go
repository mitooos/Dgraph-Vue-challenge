package test_utils

import (
	"backend/models"
	"time"
)

func randomProduct() *models.Product {
	return models.NewProduct(RandomString(8), RandomString(40), RandomInt(3000), time.Now())
}

func RandomSliceOfProducts(size int) ([]*models.Product, map[string]*models.Product) {
	products := make([]*models.Product, size)
	productsMap := make(map[string]*models.Product)
	for i := 0; i < size; i++ {
		randomProduct := randomProduct()
		products[i] = randomProduct
		productsMap[randomProduct.Id] = randomProduct
	}

	return products, productsMap
}

func MapOfProductsFromSlice(products []*models.Product) map[string]*models.Product {
	productsMap := make(map[string]*models.Product)
	for _, product := range products {
		productsMap[product.Id] = product
	}
	return productsMap
}
