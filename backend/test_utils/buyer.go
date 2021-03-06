package test_utils

import (
	"backend/models"
	"time"
)

func randomBuyer() *models.Buyer {
	return models.NewBuyer(RandomString(8), RandomString(20), RandomInt(100), time.Now())
}

func RandomSliceOfBuyers(size int) ([]*models.Buyer, map[string]*models.Buyer) {
	buyers := make([]*models.Buyer, size)
	buyersMap := make(map[string]*models.Buyer)
	for i := 0; i < size; i++ {
		randomBuyer := randomBuyer()
		buyers[i] = randomBuyer
		buyersMap[randomBuyer.Id] = randomBuyer
	}

	return buyers, buyersMap
}

func MapOfBuyersFromSlice(buyers []*models.Buyer) map[string]*models.Buyer {
	buyersMap := make(map[string]*models.Buyer)
	for _, buyer := range buyers {
		buyersMap[buyer.Id] = buyer
	}
	return buyersMap
}
