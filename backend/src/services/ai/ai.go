package ai

import (
	"context"

	"github.com/tmc/langchaingo/chains"
)

type AiService struct {
	chains []chains.Chain
}

func NewAiService() (*AiService, error) {
	return &AiService{
		chains: []chains.Chain{},
	}, nil
}

func (svc *AiService) AddChain(chain chains.Chain) {
	svc.chains = append(svc.chains, chain)
}

func (svc *AiService) Run(question string) (*string, error) {
	simpleSeqChain, err := chains.NewSimpleSequentialChain(svc.chains)
	if err != nil {
		return nil, err
	}
	answer, err := chains.Run(context.Background(), simpleSeqChain, question)
	if err != nil {
		return nil, err
	}
	return &answer, nil
}
