package main

import (
	"backend/models"
	"log"
	"net/http"
)

func init(){
	err := models.LoadSchemas()
	if(err != nil){
		log.Fatal(err)
	}
}

func main(){
	router := models.Router()
	log.Println("listening on port 5000")

	http.ListenAndServe(":5000", router)
}