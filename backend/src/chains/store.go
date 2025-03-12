package chain_store

import (
	"github.com/ivolejon/pivo/repositories"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/prompts"
)

type ChainStore struct {
	vectorStore *repositories.VectorStore
}

func NewChainStore(vstore *repositories.VectorStore) *ChainStore {
	return &ChainStore{
		vectorStore: vstore,
	}
}

func (s *ChainStore) GetBaseDocumentChain(llm llms.Model) chains.Chain {
	prompt := prompts.NewPromptTemplate(
		getBaseDocumentChainSystemPrompt(),
		[]string{"input_documents", "question"},
	)
	combineChain := chains.NewStuffDocuments(chains.NewLLMChain(llm, prompt))
	retriever := s.vectorStore.Retriver(2)
	chain := chains.NewRetrievalQA(combineChain, retriever)

	return chain
}

func (s *ChainStore) GetFormatAsDocumentChain(llm llms.Model) chains.Chain {
	prompt := prompts.NewPromptTemplate(
		getFormatAsDocumentChain(),
		[]string{"output"},
	)
	chain := chains.NewLLMChain(llm, prompt)

	return chain
}
