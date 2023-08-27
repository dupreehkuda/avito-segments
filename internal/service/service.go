package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

//go:generate mockgen -source=service.go -destination=mock_test.go -package=service_test

type UserRepository interface {
	SetSegments(ctx context.Context, segments *models.UserSetRequest) error
	DeleteSegments(ctx context.Context, segments *models.UserDeleteRequest) error
	GetSegments(ctx context.Context, userID string) (*models.UserResponse, error)

	GetReportData(ctx context.Context, year, month int) ([]models.ReportRow, error)
}

type SegmentRepository interface {
	Add(ctx context.Context, segment *models.Segment) error
	Delete(ctx context.Context, slug string) error
	Get(ctx context.Context, slug string) (*models.Segment, error)
	Count(ctx context.Context, slugs []string) (int, error)
}

// Service provides service's business-logic.
type Service struct {
	userRepo    UserRepository
	segmentRepo SegmentRepository
	logger      *zap.Logger
}

// New creates new instance of service.
func New(userRepo UserRepository, segmentRepo SegmentRepository, logger *zap.Logger) *Service {
	return &Service{
		userRepo:    userRepo,
		segmentRepo: segmentRepo,
		logger:      logger,
	}
}
