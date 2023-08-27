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
		return echo.NewHTTPError(http.StatusInternalServerError, "internal Server Error")
	}

	err = easyjson.Unmarshal(body, &req)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "internal Server Error")
	}

	if err = UUIDCheck(req.UserID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "userID should be a valid UUID")
	}

	if req.Segments == nil || len(req.Segments) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no segments provided")
	}

	if err = h.service.UserSetSegments(c.Request().Context(), &req); err != nil {
		switch {
		case errors.Is(err, errs.ErrSegmentsNotFound):
			return echo.NewHTTPError(http.StatusBadRequest, "segment(s) not found")
		case errors.Is(err, errs.ErrAlreadyExpired):
			return echo.NewHTTPError(http.StatusBadRequest, "segment operation expired")
		case errors.Is(err, errs.ErrInvalidSegmentSlug):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid slug naming")
		default:
			h.logger.Error("Error occurred setting segments", zap.Error(err))
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) UserDeleteSegments(c echo.Context) error {
	var req models.UserDeleteRequest

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

	if err = UUIDCheck(req.UserID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "userID should be a valid UUID")
	}

	if req.Slugs == nil || len(req.Slugs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no segments provided")
	}

	if err = h.service.UserDeleteSegments(c.Request().Context(), &req); err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidSegmentSlug):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid slug naming")
		default:
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handlers) UserGetSegments(c echo.Context) error {
	id := c.Param("id")

	if err := UUIDCheck(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "userID should be a valid UUID")
	}

	resp, err := h.service.UserGetSegments(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrSegmentsNotFound):
			return c.NoContent(http.StatusNoContent)
		case errors.Is(err, errs.ErrUserNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		default:
			return err
		}
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
			return err
		}
	}

	return nil
}
