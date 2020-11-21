package models

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dgraph-io/dgo/v200/protos/api"
)


type Buyer struct{
	Id string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Age int `json:"age,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

func NewBuyer(id string, name string, age int) *Buyer{
	return &Buyer{
		Id: id,
		Name: name,
		Age: age,
		DType: []string{"Buyer"},
	}
}

const buyerSchema = `
id: string @index(exact) .
name: string .
age: int .
`

func InsertManyBuyers(buyers []*Buyer) error {
	c, err := NewClient()
	if err != nil{
		log.Panic(err)
		return err
	}
	txn := c.NewTxn()
	defer txn.Discard(context.Background())

	out, err := json.Marshal(buyers)
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