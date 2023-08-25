package service

import (
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=service

type Repository interface {
	AddSegment(segment models.Segment) error
	DeleteSegment(tag string) error
}

// Service provides service's business-logic
type Service struct {
	repository Repository
	logger     *zap.Logger
}

// New creates new instance of service
func New(Repository Repository, logger *zap.Logger) *Service {
	return &Service{
		repository: Repository,
		logger:     logger,
	}
}
