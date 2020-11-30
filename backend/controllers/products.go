package controllers

import (
	"backend/models"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"time"
)

func insertProducts(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.Panic(err)
		respondWithError(w, 400, "Unable to get file")
		return
	}

	reader := csv.NewReader(file)
	reader.Comma = '\''
	record, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
		respondWithError(w, 400, "Unable to read csv")
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		log.Panic(err)
		respondWithError(w, 400, "Cannot parse date field")
		return
	}

	var price int
	products := make([]*models.Product, len(record))
	for i, line := range record {
		price, _ = strconv.Atoi(line[2])
		product := models.NewProduct(line[0], line[1], price, date)
		products[i] = product
	}

	err = models.InsertManyProducts(products)
	if err != nil {
		respondWithError(w, 500, "Unable to load data")
		return
	}

	respondWithJSON(w, 200, products)
}
