package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h Handlers) ReportCreate(c echo.Context) error {
	// var req models.ReportRequest
	//
	// body, err := io.ReadAll(c.Request().Body)
	// if err != nil {
	//	h.logger.Error("Unable to read body", zap.Error(err))
	//	return echo.NewHTTPError(http.StatusInternalServerError, "internal Server Error")
	// }
	//
	// err = easyjson.Unmarshal(body, &req)
	// if err != nil {
	//	h.logger.Error("Unable to decode JSON", zap.Error(err))
	//	return echo.NewHTTPError(http.StatusInternalServerError, "internal Server Error")
	// }
	//
	// fileName, err := h.service.CreateReport(c.Request().Context(), req.Year, req.Month)
	// if err != nil {
	//
	// }

	h.logger.Info("test info, getting host", zap.String("host", c.Request().Host))

	return nil
}

func (h Handlers) ReportGet(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
