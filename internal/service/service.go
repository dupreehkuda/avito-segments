package service

import "go.uber.org/zap"

type Repository interface {
}

// Service provides service's business-logic
type Service struct {
	repository Repository
	logger     *zap.Logger
}

// New creates new instance of processor
func New(Repository Repository, logger *zap.Logger) *Service {
	return &Service{
		repository: Repository,
		logger:     logger,
	}
}
