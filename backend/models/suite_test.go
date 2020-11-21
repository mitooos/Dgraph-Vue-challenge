package models_test

import (
	"backend/models"
	"math/rand"
	"os"
	"testing"
	"time"
)


func TestMain(m *testing.M){
	models.LoadSchemas()

	code := m.Run()

	os.Exit(code)
}


func randomString(length int) string {

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func randomInt(max int) int {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	
	return seededRand.Intn(max)
}