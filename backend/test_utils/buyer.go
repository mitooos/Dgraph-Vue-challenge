package test_utils

import "backend/models"

func randomBuyer() *models.Buyer {
	return models.NewBuyer(randomString(8), randomString(20), randomInt(100))
}

func RandomSliceOfBuyers(size int)([]*models.Buyer,map[string]*models.Buyer) {
	buyers := make([]*models.Buyer, size)
	buyersMap := make(map[string]*models.Buyer)
	for i:=0; i < size; i++{
		randomBuyer := randomBuyer()
		buyers[i] = randomBuyer
		buyersMap[randomBuyer.Id] = randomBuyer
	}

	return buyers, buyersMap
}

func MapOfBuyersFromSlice(buyers []*models.Buyer) map[string]*models.Buyer {
	buyersMap := make(map[string]*models.Buyer)
	for _, buyer := range buyers{
		buyersMap[buyer.Id] = buyer
	}
	return buyersMap
}