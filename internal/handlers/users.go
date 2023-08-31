package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	errs "github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (h Handlers) UserSetSegments(c echo.Context) error {
	var req models.UserSetRequest

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

	if err = UUIDCheck(req.UserID); err != nil {
		return h.ErrorHandler(err)
	}

	if req.Segments == nil || len(req.Segments) == 0 {
		return h.ErrorHandler(errs.ErrNoSegmentsProvided)
	}

	if err = h.service.UserSetSegments(c.Request().Context(), &req); err != nil {
		return h.ErrorHandler(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) UserDeleteSegments(c echo.Context) error {
	var req models.UserDeleteRequest

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

	if err = UUIDCheck(req.UserID); err != nil {
		return h.ErrorHandler(err)
	}

	if req.Slugs == nil || len(req.Slugs) == 0 {
		return h.ErrorHandler(errs.ErrNoSegmentsProvided)
	}

	if err = h.service.UserDeleteSegments(c.Request().Context(), &req); err != nil {
		return h.ErrorHandler(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) UserGetSegments(c echo.Context) error {
	id := c.Param("id")

	if err := UUIDCheck(id); err != nil {
		return h.ErrorHandler(err)
	}

	resp, err := h.service.UserGetSegments(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrSegmentsNotFound) {
			return c.NoContent(http.StatusNoContent)
		}

		return h.ErrorHandler(err)
	}

	return c.JSON(http.StatusOK, resp)
}

// UUIDCheck checks all request ids to be sure they're not empty and correct uuids.
func UUIDCheck(uuids ...string) error {
	for _, id := range uuids {
		if id == "" {
			return errs.ErrInvalidUserID
		}

		_, err := uuid.Parse(id)
		if err != nil {
			return errs.ErrInvalidUserID
		}
	}

	return nil
}
