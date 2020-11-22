package models_test

import (
	"backend/models"
	"os"
	"testing"
)


func TestMain(m *testing.M){
	models.LoadSchemas()

	code := m.Run()

	os.Exit(code)
}