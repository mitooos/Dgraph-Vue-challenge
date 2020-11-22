package controllers_test

import (
	"backend/controllers"
	"backend/models"
	"os"
	"testing"

	"github.com/go-chi/chi"
)

var router *chi.Mux


func TestMain(m *testing.M){
	models.LoadSchemas()
	router = controllers.Router()

	code := m.Run()

	os.Exit(code)
}