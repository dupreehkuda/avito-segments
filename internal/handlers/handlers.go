package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	errs "github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

//go:generate mockgen -source=handlers.go -destination=mock_test.go -package=handlers_test

// Service is an interface for business-logic.
type Service interface {
	SegmentAdd(ctx context.Context, segment *models.Segment) error
	SegmentDelete(ctx context.Context, slug string) error
	CreateReport(ctx context.Context, year, month int) (string, error)

	UserSetSegments(ctx context.Context, segments *models.UserSetRequest) error
	UserDeleteSegments(ctx context.Context, segments *models.UserDeleteRequest) error
	UserGetSegments(ctx context.Context, userID string) (*models.UserResponse, error)
}

// Handlers provide access to service.
type Handlers struct {
	service Service
	logger  *zap.Logger
}

// New creates new instance of handlers.
func New(service Service, logger *zap.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

func (h Handlers) ErrorHandler(err error) *echo.HTTPError {
	switch {
	case errors.Is(err, errs.ErrInvalidPeriod):
		return echo.NewHTTPError(http.StatusBadRequest, "invalid time period provided")
	case errors.Is(err, errs.ErrDataNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "no data for report")
	case errors.Is(err, errs.ErrInvalidSegmentSlug):
		return echo.NewHTTPError(http.StatusBadRequest, "invalid slug naming")
	case errors.Is(err, errs.ErrInvalidUserID):
		return echo.NewHTTPError(http.StatusBadRequest, "invalid userID")
	case errors.Is(err, errs.ErrNoSegmentsProvided):
		return echo.NewHTTPError(http.StatusBadRequest, "no segments provided")
	case errors.Is(err, errs.ErrReportNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "requested report not found")
	case errors.Is(err, errs.ErrSegmentsNotFound):
		return echo.NewHTTPError(http.StatusBadRequest, "segment(s) not found")
	case errors.Is(err, errs.ErrAlreadyExpired):
		return echo.NewHTTPError(http.StatusBadRequest, "segment operation expired")
	case errors.Is(err, errs.ErrUserNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	case errors.Is(err, errs.ErrSegmentNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "slug not found")
	case errors.Is(err, errs.ErrAlreadyDeleted):
		return echo.NewHTTPError(http.StatusGone, "slug has been already deleted")
	default:
		h.logger.Error("Error occurred creating report", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
}
