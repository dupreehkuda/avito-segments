package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	errs "github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (h Handlers) ReportCreate(c echo.Context) error {
	var req models.ReportRequest

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

	potentialFileName := fmt.Sprintf("%v_%v_report.csv", req.Month, req.Year)
	if _, err = os.Stat("reports/" + potentialFileName); !errors.Is(err, os.ErrNotExist) {
		return c.JSON(http.StatusOK, models.ReportResponse{Link: c.Request().Host + "/api/v1/report/" + potentialFileName})
	}

	fileName, err := h.service.CreateReport(c.Request().Context(), req.Year, req.Month)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidPeriod):
			return echo.NewHTTPError(http.StatusBadRequest, "invalid time period provided")
		case errors.Is(err, errs.ErrDataNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "no data for report")
		default:
			h.logger.Error("Error occurred creating report", zap.Error(err))
			return err
		}
	}

	resp := models.ReportResponse{Link: c.Request().Host + "/api/v1/report/" + fileName}

	return c.JSON(http.StatusOK, resp)
}

func (h Handlers) ReportGet(c echo.Context) error {
	file := c.Param("file")

	if _, err := os.Stat("reports/" + file); errors.Is(err, os.ErrNotExist) {
		return echo.NewHTTPError(http.StatusBadRequest, "report doesn't exist")
	}

	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment;filename="+file)

	return c.File("reports/" + file)
}
