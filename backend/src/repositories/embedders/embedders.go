package embedders

import (
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/ztrue/tracerr"
)

func GetEmbedderNomicEmbedTextModel() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("nomic-embed-text"))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return embedder, nil
}

func GetEmbedderLlama2_3Model() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return embedder, nil
}

func GetEmbedderGemma3_27b() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("gemma3:27b"))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return embedder, nil
}

func GetEmbedderBgeLarge() (*embeddings.EmbedderImpl, error) {
	llm, err := ollama.New(ollama.WithModel("bge-large"))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return embedder, nil
}
