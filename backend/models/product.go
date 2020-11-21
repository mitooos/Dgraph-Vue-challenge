package models

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v200/protos/api"
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

const productSchema = `
id: string @index(exact) .
name: string .
price: int .
`

func InsertManyProducts(products []*Product) error{
	c, err := NewClient()
	if err != nil{
		log.Panic(err)
		return err
	}
	txn := c.NewTxn()
	defer txn.Discard(context.Background())

	out, err := json.Marshal(products)
	if err != nil{
		log.Panic(err)
		return err
	}

	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: out})
	if err != nil{
		log.Panic(err)
		return err
	}else{
		err = txn.Commit(context.Background())
		if err != nil{
			log.Panic(err)
			return err
		}
	}
	return nil
}