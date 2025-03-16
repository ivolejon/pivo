package embedders

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func GetEmbedderNomicEmbedTextModel() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("nomic-embed-text"))
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	return embedder, nil
}

func GetEmbedderLlama2_3Model() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return nil, err
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, err
	}
	return embedder, nil
}
