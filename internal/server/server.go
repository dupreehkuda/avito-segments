package server

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
)

type Api struct {
	handlers Handlers
	config   *config.Config
	logger   *zap.Logger
}

func RegisterServer(lc fx.Lifecycle, handlers Handlers, config *config.Config, logger *zap.Logger) *Api {
	api := &Api{
		handlers: handlers,
		config:   config,
		logger:   logger,
	}

	serv := &http.Server{Addr: config.Conn.Port, Handler: api.handler(logger)}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := serv.ListenAndServe()
				if err != nil {
					logger.Fatal("Cant start server", zap.Error(err))
				}
			}()

			logger.Info("Server started", zap.String("port", config.Conn.Port))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := serv.Shutdown(ctx)
			if err != nil {
				logger.Fatal("Error shutting down", zap.Error(err))
				return err
			}

			logger.Info("Server shut down", zap.String("port", config.Conn.Port))

			return nil
		},
	})

	return api
}
