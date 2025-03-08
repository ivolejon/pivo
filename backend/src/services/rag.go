package services

type RAGService struct{}

func NewRAGService() *RAGService {
	return &RAGService{}
}

func (r *RAGService) AddDocument() string {}

func (r *RAGService) SimilaritySearch() string {}

func (r *RAGService) UpdateStore() string {}
