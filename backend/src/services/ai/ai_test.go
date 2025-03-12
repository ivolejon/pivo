package ai_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

func getOllama() (*ollama.LLM, error) {
	llm, err := ollama.New(ollama.WithModel("llama3.2"))
	if err != nil {
		return nil, err
	}
	return llm, nil
}

func TestSequentialWithTwoChains(t *testing.T) {
	t.Parallel()
	llm, err := getOllama()
	require.NoError(t, err)
	steps := []chains.Chain{
		chains.NewLLMChain(llm, prompts.NewPromptTemplate("{{.input}}", []string{"input"})),
		chains.NewLLMChain(llm, prompts.NewPromptTemplate(`Take the text a contruct a json string in format {"content: text"}: {{.output}}`, []string{"output"})),
	}
	simpleSeqChain, err := chains.NewSimpleSequentialChain(steps)
	require.NoError(t, err)

	res, err := chains.Run(context.Background(), simpleSeqChain, "Write a short sentence in english")
	require.NoError(t, err)
	require.Contains(t, res, `{"content"`)
}

func TestSequentialWithOneChain(t *testing.T) {
	t.Parallel()
	llm, err := getOllama()
	require.NoError(t, err)
	steps := []chains.Chain{
		chains.NewLLMChain(llm, prompts.NewPromptTemplate("{{.input}}", []string{"input"})),
	}
	simpleSeqChain, err := chains.NewSimpleSequentialChain(steps)
	require.NoError(t, err)

	res, err := chains.Run(context.Background(), simpleSeqChain,
		`Write a short sentence in english, contruct a json string in format {"content: text"}`)
	require.NoError(t, err)
	require.Contains(t, res, `{"content"`)
}
