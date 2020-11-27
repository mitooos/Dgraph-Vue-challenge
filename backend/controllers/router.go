package controllers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Router()*chi.Mux{
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/products", insertProducts)
	r.Post("/buyers", insertBuyers)
	r.Post("/transactions", insertTransactions)

	r.Get("/buyers", getBuyers)
	r.Get("/buyers/{buyerId}", getBuyer)


	return r
}