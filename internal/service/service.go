package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=service_test

type Repository interface {
	SegmentAdd(ctx context.Context, segment *models.Segment) error
	SegmentDelete(ctx context.Context, slug string) error
	SegmentGet(ctx context.Context, slug string) (*models.Segment, error)
	SegmentCount(ctx context.Context, slugs []string) (int, error)

	UserSetSegments(ctx context.Context, segments *models.UserSetRequest) error
	UserDeleteSegments(ctx context.Context, segments *models.UserDeleteRequest) error
	UserGetSegments(ctx context.Context, userID string) (*models.UserResponse, error)
}

// Service provides service's business-logic.
type Service struct {
	repository Repository
	logger     *zap.Logger
}

// New creates new instance of service.
func New(repository Repository, logger *zap.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}
