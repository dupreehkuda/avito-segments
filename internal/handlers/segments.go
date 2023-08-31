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
		return h.ErrorHandler(err)
	}

	err = easyjson.Unmarshal(body, &req)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		return h.ErrorHandler(err)
	}

	if err = h.service.SegmentAdd(c.Request().Context(), &req); err != nil {
		if errors.Is(err, errs.ErrDuplicateSegment) {
			return c.NoContent(http.StatusOK)
		}

		return h.ErrorHandler(err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h Handlers) SegmentDelete(c echo.Context) error {
	slug := c.Param("slug")

	if err := h.service.SegmentDelete(c.Request().Context(), slug); err != nil {
		return h.ErrorHandler(err)
	}

	return c.NoContent(http.StatusOK)
}
