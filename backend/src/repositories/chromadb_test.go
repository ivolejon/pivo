package repositories_test

import (
	"testing"

	"github.com/ivolejon/pivo/repositories"
)

func TestNewChromaDB(t *testing.T) {
	db, err := repositories.NewChromaDB()
	if err != nil {
		t.Errorf("Error creating ChromaDB: %v", err)
	}
	if db == nil {
		t.Errorf("ChromaDB is nil")
	}
}
