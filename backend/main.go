package main

import (
	"backend/controllers"
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
	router := controllers.Router()
	log.Println("listening on port 5000")

	log.Fatal(http.ListenAndServe(":5000", router))
}