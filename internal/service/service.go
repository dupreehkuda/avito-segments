package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=service

type Repository interface {
	AddSegment(ctx context.Context, segment models.Segment) error
	DeleteSegment(ctx context.Context, tag string) error
	GetSegment(ctx context.Context, tag string) (*models.Segment, error)
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
