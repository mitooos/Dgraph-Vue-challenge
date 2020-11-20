package models

import (
	"context"
	"log"

	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"google.golang.org/grpc"
)

func newClient() (*dgo.Dgraph, error){
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil{
		log.Print("Error creating the client")
		log.Panic(err)
		return nil, err
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	), nil
}

func setupDb(c *dgo.Dgraph) error{
	err := c.Alter(context.Background(), &api.Operation{DropOp: api.Operation_ALL})
	if err != nil{
		log.Println("Error droping all data")
		log.Panic(err)
		return err
	}
	
	err = c.Alter(context.Background(), &api.Operation{
		Schema: productSchema,
	})
	if err != nil{
		log.Print("Error loading product schema")
		log.Panic(err)
		return err
	}
	return nil
}

func LoadSchemas() error{
	c, err := newClient()
	if err != nil{
		return err
	}
	
	return setupDb(c)
}

