package controllers

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func insertBuyers(w http.ResponseWriter, r *http.Request){
	file, _, err := r.FormFile("file")
	defer file.Close()
	if(err != nil){
		log.Panic(err)
		respondWithError(w, 400, "Unable to get file")
		return
	}

	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err !=nil{
		log.Panic(err)
		respondWithError(w, 400, "Unable to read json")
		return
	}

	var buyers []*models.Buyer

	if err := json.Unmarshal(buffer.Bytes(), &buyers); err != nil {
		log.Panic(err)
		respondWithError(w, 400, "Unable to read json")
		return
	}

	for i, _ := range buyers{
		buyers[i].DType = []string{"Buyer"}
	}

	err = models.InsertManyBuyers(buyers)
	if err != nil{
		log.Panic(err)
		respondWithError(w, 500, "Unable to load data")
		return
	}

	respondWithJSON(w, 200, buyers)
}

func getBuyers(w http.ResponseWriter, r *http.Request){
	buyers, err := models.GetBuyers()
	if err != nil{
		respondWithError(w, 500, "Unable to get buyers")
		return
	}

	respondWithJSON(w, 200, buyers)
}