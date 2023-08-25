package handlers

import (
	"context"

	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

//go:generate mockgen -source=handlers.go -destination=mock_test.go -package=handlers

// Service is an interface for business-logic
type Service interface {
	SegmentAdd(ctx context.Context, segment models.Segment) error
	SegmentDelete(ctx context.Context, tag string) error
}

// Handlers provide access to service
type Handlers struct {
	service Service
	logger  *zap.Logger
}

// New creates new instance of handlers
func New(service Service, logger *zap.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}
