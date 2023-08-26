package server

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
)

type API struct {
	handlers Handlers
	config   *config.Config
	logger   *zap.Logger
}

func RegisterServer(lc fx.Lifecycle, handlers Handlers, config *config.Config, logger *zap.Logger) *API {
	api := &API{
		handlers: handlers,
		config:   config,
		logger:   logger,
	}

	serv := api.handler(logger)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := serv.Start(config.Conn.Port)
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
