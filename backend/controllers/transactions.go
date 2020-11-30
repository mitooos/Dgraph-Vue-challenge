package controllers

import (
	"backend/models"
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func insertTransactions(w http.ResponseWriter, r *http.Request) {
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

	transactionsStringSlice := strings.Split(buffer.String(), "#")

	date, err := time.Parse("2006-01-02", r.FormValue("date"))
	if err != nil {
		log.Print(err)
		respondWithError(w, 400, "Cannot parse date field")
		return
	}

	for _, transactionString := range transactionsStringSlice {
		if len(transactionString) != 0 {
			// TODO: make concurrent
			transactionStringSlice := strings.Split(transactionString, "\000")
			if len(transactionStringSlice) < 5 {
				respondWithError(w, 400, "Data does not match format")
				return
			}
			id := transactionStringSlice[0]
			buyerId := transactionStringSlice[1]
			ip := transactionStringSlice[2]
			device := transactionStringSlice[3]

			productsSubString := transactionStringSlice[4]
			products := strings.ReplaceAll(productsSubString[1:len(productsSubString)-1], ",", " ")

			if err := models.InsertTransaction(id, date, buyerId, ip, device, products); err != nil {
				log.Fatal(err)
			}
		}
	}

	respondWithJSON(w, 200, map[string]string{"message": "success"})
}
