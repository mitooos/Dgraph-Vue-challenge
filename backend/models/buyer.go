package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Buyer struct {
	Id           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Age          int           `json:"age,omitempty"`
	Date         time.Time     `json:"date,omiempty"`
	DType        []string      `json:"dgraph.type,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

func NewBuyer(id string, name string, age int, date time.Time) *Buyer {
	return &Buyer{
		Id:    id,
		Name:  name,
		Age:   age,
		Date:  date,
		DType: []string{"Buyer"},
	}
}

func InsertManyBuyers(buyers []*Buyer) error {
	out, err := json.Marshal(buyers)
	if err != nil {
		log.Panic(err)
		return err
	}

	return ExecuteMutation(out)
}

func GetBuyers(first string, offset string) ([]*Buyer, error) {
	query := fmt.Sprintf(`
		{
			Buyers(func: type(Buyer), first: %s, offset: %s){
				id
				name
				age
				date
				dgraph.type
			}
		}
	`, first, offset)

	rep, err := ExecuteQuery(query)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	var buyers struct {
		Buyers []*Buyer
	}

	if err := json.Unmarshal(rep.Json, &buyers); err != nil {
		log.Panic(err)
		return nil, err
	}
	return buyers.Buyers, nil
}

type BuyerDetailResponse struct {
	Buyer               []*Buyer
	BuyersWithSameIP    []*Buyer
	RecommendedProducts []*struct {
		Id    string `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Count int    `json:"count,omitempty"`
	}
}

func GetBuyer(id string) (*BuyerDetailResponse, error) {
	query := fmt.Sprintf(`
	{
		var(func: eq(id, "%s")){
    		buyerUid as uid
			transactions{
				buyerIps as ip
        		products{
          			transactionsWithSameProducts as ~products
        		}
			}
		}
  
		Buyer(func: uid(buyerUid)){
			id
			name
			age
			date
			dgraph.type
			transactions{
				id
				date
				ip
				device
				dgraph.type
				products{
					id
					name
					price
					dgraph.type
				}
			}
		}

		var(func: eq(ip, val(buyerIps))){
			buyer @filter(NOT uid(buyerUid)){
				buyersSameIp as uid
			}
		}
  
		BuyersWithSameIP(func: uid(buyersSameIp)){
			id
			name
			age
			date
			dgraph.type
		}
  
  		recomProds as var(func: type(Product)){ 
			prodsCount as count(~products) @filter(uid(transactionsWithSameProducts))
  		}
    
  		RecommendedProducts(func: uid(recomProds), orderdesc: val(prodsCount), first: 10){
			id
			name
    		count: val(prodsCount)
  		}
  
	}`, id)

	rep, err := ExecuteQuery(query)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	var response BuyerDetailResponse

	if err := json.Unmarshal(rep.Json, &response); err != nil {
		log.Panic(err)
		return nil, err
	}

	return &response, nil
}
