package main

import (
	"backend/controllers"
	"backend/models"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func init(){
	err := models.LoadSchemas()
	if(err != nil){
		log.Fatal(err)
	}
}

func main(){
	router := router()
	log.Println("listening on port 5000")

	http.ListenAndServe(":5000", router)
}


func router()*chi.Mux{
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/products", controllers.InsertProducts)


	return r
}