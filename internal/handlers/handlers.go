package handlers

import "go.uber.org/zap"

// Service is an interface for business-logic
type Service interface {
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
