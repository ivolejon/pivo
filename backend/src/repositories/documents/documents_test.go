package documents_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/repositories/documents"
	"github.com/stretchr/testify/require"
)

func TestCanInitlizeRepo(t *testing.T) {
	_, err := documents.NewDocumentsRepository()
	require.NoError(t, err)
}

func TestAddDocment(t *testing.T) {
	repo, err := documents.NewDocumentsRepository()
	require.NoError(t, err)

	ID := uuid.New()
	projectID := uuid.New()
	title := "New doc"
	now := time.Now()
	embeddID1 := uuid.New()
	embeddID2 := uuid.New()
	embeddingsIds := []uuid.UUID{embeddID1, embeddID2}

	params := documents.AddDocumentParams{
		ID:            ID,
		EmbeddingsIds: embeddingsIds,
		Filename:      "IvoPivoDoc.pdf",
		Title:         &title,
		ProjectID:     projectID,
		CreatedAt:     now,
	}

	newDoc, err := repo.AddDocument(params)
	require.NoError(t, err)
	if !now.Equal(newDoc.CreatedAt) {
		t.Errorf("Expected %v, got %v", now, newDoc.CreatedAt)
	}
	require.Equal(t, title, *newDoc.Title)
}
