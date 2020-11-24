package models

import (
	"encoding/json"
	"log"
)

type Product struct{
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Price int `json:"price,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

func NewProduct(id string, name string, price int)*Product{
	return &Product{
		Id: id,
		Name: name,
		Price: price,
		DType: []string{"Product"},
	}
}

func InsertManyProducts(products []*Product) error {
	out, err := json.Marshal(products)
	if err != nil{
		log.Panic(err)
		return err
	}

	return ExecuteMutation(out)
}