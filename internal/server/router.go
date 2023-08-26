package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Handlers interface {
	SegmentAdd(c echo.Context) error
	SegmentDelete(c echo.Context) error
}

func (a *API) handler(logger *zap.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.String("remote", v.RemoteIP),
				zap.Duration("latency", v.Latency),
			)

			return nil
		},
	}))

	api := e.Group("/api")
	v1 := api.Group("/v1")

	segment := v1.Group("/segment")

	segment.POST("", a.handlers.SegmentAdd)
	segment.DELETE("/:tag", a.handlers.SegmentDelete)

	user := v1.Group("/user")

	user.GET("/", nil)
	user.POST("/", nil)
	user.DELETE("/", nil)

	return e
}
