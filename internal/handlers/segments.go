package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	errs "github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (h Handlers) SegmentAdd(c echo.Context) error {
	var req models.Segment

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("Unable to read body", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "internal Server Error")
	}

	err = easyjson.Unmarshal(body, &req)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "internal Server Error")
	}

	if err = h.service.SegmentAdd(c.Request().Context(), &req); err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateSegment):
			return c.NoContent(http.StatusOK)
		case errors.Is(err, errs.ErrInvalidSegmentSlug):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid slug naming")
		default:
			h.logger.Error("Error occurred adding segment", zap.Error(err))
			return err
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (h Handlers) SegmentDelete(c echo.Context) error {
	slug := c.Param("slug")

	if err := h.service.SegmentDelete(c.Request().Context(), slug); err != nil {
		switch {
		case errors.Is(err, errs.ErrSegmentNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "slug not found")
		case errors.Is(err, errs.ErrAlreadyDeleted):
			return echo.NewHTTPError(http.StatusGone, "slug has been already deleted")
		case errors.Is(err, errs.ErrInvalidSegmentSlug):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid slug naming")
		default:
			h.logger.Error("Error occurred deleting segment", zap.Error(err))
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
