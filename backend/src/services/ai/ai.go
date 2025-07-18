package ai

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/chains"
	"github.com/ztrue/tracerr"
)

type AiService struct {
	chains     []chains.Chain
	streamChan chan AiStreamingEvent
}

type AiStreamingEvent struct {
	Chunk  []byte
	Error  error
	Status string // e.g., "streaming", "completed", "error"
}

func NewAiService() (*AiService, error) {
	return &AiService{
		chains:     []chains.Chain{},
		streamChan: make(chan AiStreamingEvent),
	}, nil
}

func (svc *AiService) GetStreamChan() chan AiStreamingEvent {
	return svc.streamChan
}

func (svc *AiService) AddChain(chain chains.Chain) {
	svc.chains = append(svc.chains, chain)
}

func (svc *AiService) Run(question string) error {
	simpleSeqChain, err := chains.NewSimpleSequentialChain(svc.chains)
	if err != nil {
		return tracerr.Wrap(err)
	}

	streamingFunc := func(ctx context.Context, chunk []byte) error {
		log.Printf("Received chunk: %s", chunk)
		svc.streamChan <- AiStreamingEvent{Chunk: chunk, Status: "streaming"}
		return nil
	}

	streamingOption := chains.WithStreamingFunc(streamingFunc)

	options := []chains.ChainCallOption{
		streamingOption,
	}

	answer, err := chains.Run(context.Background(), simpleSeqChain, question, options...)
	if err != nil {
		svc.streamChan <- AiStreamingEvent{Error: err, Status: "error"}
	}
	svc.streamChan <- AiStreamingEvent{Chunk: []byte(answer), Status: "completed"}
	return nil
}
