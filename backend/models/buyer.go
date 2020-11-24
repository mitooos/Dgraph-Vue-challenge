package models

import (
	"encoding/json"
	"log"
)


type Buyer struct{
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age int `json:"age,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

func NewBuyer(id string, name string, age int) *Buyer{
	return &Buyer{
		Id: id,
		Name: name,
		Age: age,
		DType: []string{"Buyer"},
	}
}

func InsertManyBuyers(buyers []*Buyer) error {
	out, err := json.Marshal(buyers)
	if err != nil{
		log.Panic(err)
		return err
	}
	
	return ExecuteMutation(out)
}