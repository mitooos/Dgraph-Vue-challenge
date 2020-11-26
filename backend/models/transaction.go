package models

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)


type Transaction struct{
	Id string `json:"id,omitempty"`
	Date time.Time `json:"date,omiempty"`
	Buyer *Buyer `json:"buyer,omiempty"`
	Ip string `json:"ip,omitempty"`
	Device string `json:"device,omitempty"`
	Products []*Product `json:"products,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
}

func InsertTransaction(id string, date time.Time, buyerId string, ip string, device string, productIds string) error {

	buyerUidQuery := fmt.Sprintf(`
		{
			Get(func: eq(id, "%s")){
				uid
			}
		}
	`, buyerId)
	res, err := ExecuteQuery(buyerUidQuery)
	if err != nil{
		panic(err)
		return err
	}

	buyerUidDecode := struct{
		Get []struct{
			Uid *string
		}
	}{}
	if err := json.Unmarshal(res.GetJson(), &buyerUidDecode); err != nil{
		panic(err)
		return err
	}


	buyerUid := *buyerUidDecode.Get[0].Uid
	

	productsUidQuery := fmt.Sprintf(`
		{
			Get(func: anyofterms(id, "%s")){
				uid
			}
		}
	`, productIds)
	res, err = ExecuteQuery(productsUidQuery)
	if err != nil{
		panic(err)
		return err
	}

	productsUidDecode := struct{
		Get []struct{
			Uid *string
		}
	}{}
	if err := json.Unmarshal(res.GetJson(), &productsUidDecode); err != nil{
		panic(err)
		return err
	}

	productIdsMutation := ""
	for _, productUid := range productsUidDecode.Get{
		productIdsMutation += ("_:trans <products> <" + *productUid.Uid + "> .\n")
	}

	mutation := fmt.Sprintf(`
		_:trans <id> "%s" .
		_:trans <date> "%s" .
		_:trans <buyer>  <%s> .
		<%s> <transactions> _:trans  .
		_:trans <ip>  "%s" .
		_:trans <device>  "%s" .
		_:trans <dgraph.type>  "Transaction" .
		%s
	`, id, date.Format("2006-01-02T15:04:05"), buyerUid, buyerUid, ip, device, productIdsMutation)
	
		if err := ExecuteMutationNQuad([]byte(mutation)); err != nil{
			log.Panic(err)
			return err
		}

		return nil
}