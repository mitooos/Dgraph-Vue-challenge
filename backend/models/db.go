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
		Schema: schema,
	})
	if err != nil{
		log.Print("Error loading schemas")
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

func ExecuteQuery(query string)(*api.Response, error){
	c, err := newClient()
	if err != nil{
		return nil, err
	}
	txn := c.NewTxn()

	return txn.Query(context.Background(), query)
}

func ExecuteMutation(mutation []byte) error{
	c, err := newClient()
	if err != nil{
		log.Panic(err)
		return err
	}
	txn := c.NewTxn()
	defer txn.Discard(context.Background())
	

	_, err = txn.Mutate(context.Background(), &api.Mutation{SetJson: mutation})
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

func ExecuteMutationNQuad(mutation []byte) error{
	c, err := newClient()
	if err != nil{
		log.Panic(err)
		return err
	}
	txn := c.NewTxn()
	defer txn.Discard(context.Background())
	

	_, err = txn.Mutate(context.Background(), &api.Mutation{SetNquads: mutation, CommitNow: true})
	if err != nil{
		log.Panic(err)
		return err
	}
	return nil
}

