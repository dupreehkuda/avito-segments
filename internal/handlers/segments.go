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

func (h Handlers) AddSegment(c echo.Context) error {
	var req models.Segment

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error("Unable to read body", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if err = easyjson.Unmarshal(body, &req); err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if err := h.service.AddSegment(c.Request().Context(), req); err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateSegment):
			return c.NoContent(http.StatusOK)
		case errors.Is(err, errs.ErrInvalidSegmentTag):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid tag naming")
		default:
			return err
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (h Handlers) DeleteSegment(c echo.Context) error {
	tag := c.Param("tag")

	if err := h.service.DeleteSegment(c.Request().Context(), tag); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "Tag not found")
		case errors.Is(err, errs.ErrAlreadyDeleted):
			return echo.NewHTTPError(http.StatusGone, "Tag has been already deleted")
		case errors.Is(err, errs.ErrInvalidSegmentTag):
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid tag naming")
		default:
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
