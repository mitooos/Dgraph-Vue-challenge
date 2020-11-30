package controllers

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func insertBuyers(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		log.Print(err)
		respondWithError(w, 400, "Unable to get file")
		return
	}

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		log.Print(err)
		respondWithError(w, 400, "Unable to read json")
		return
	}

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		log.Print(err)
		respondWithError(w, 400, "Cannot parse date field")
		return
	}

	var buyers []*models.Buyer

	if err := json.Unmarshal(buffer.Bytes(), &buyers); err != nil {
		log.Print(err)
		respondWithError(w, 400, "Unable to read json")
		return
	}

	for i, _ := range buyers {
		buyers[i].Date = date
	}

	for i, _ := range buyers {
		buyers[i].DType = []string{"Buyer"}
	}

	err = models.InsertManyBuyers(buyers)
	if err != nil {
		log.Print(err)
		respondWithError(w, 500, "Unable to load data")
		return
	}

	respondWithJSON(w, 200, buyers)
}

func getBuyers(w http.ResponseWriter, r *http.Request) {
	first, ok := r.URL.Query()["first"]
	if !ok {
		first = []string{"10000"}
	}

	offset, ok := r.URL.Query()["offset"]
	if !ok {
		offset = []string{"0"}
	}

	buyers, err := models.GetBuyers(first[0], offset[0])
	if err != nil {
		respondWithError(w, 500, "Unable to get buyers")
		return
	}

	respondWithJSON(w, 200, buyers)
}

func getBuyer(w http.ResponseWriter, r *http.Request) {
	buyerID := chi.URLParam(r, "buyerId")
	resp, err := models.GetBuyer(buyerID)
	if err != nil {
		respondWithError(w, 500, "Unable to get buyer")
		return
	}

	if len(resp.Buyer) == 0 {
		respondWithError(w, 404, "Buyer not found")
		return
	}

	respondWithJSON(w, 200, resp)
}
