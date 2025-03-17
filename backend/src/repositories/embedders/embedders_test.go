package embedders_test

import (
	"testing"

	"github.com/ivolejon/pivo/repositories/embedders"
	"github.com/stretchr/testify/require"
)

func TestGetEmbedderNomicEmbedTextModel(t *testing.T) {
	_, err := embedders.GetEmbedderNomicEmbedTextModel()
	require.NoError(t, err)
}

func TestGetEmbedderLlama2_3Model(t *testing.T) {
	_, err := embedders.GetEmbedderLlama2_3Model()
	require.NoError(t, err)
}
