package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
	"github.com/dupreehkuda/avito-segments/internal/handlers"
	"github.com/dupreehkuda/avito-segments/internal/logger"
	"github.com/dupreehkuda/avito-segments/internal/repository"
	"github.com/dupreehkuda/avito-segments/internal/server"
	"github.com/dupreehkuda/avito-segments/internal/service"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(config.New),
		fx.Provide(logger.New),
		fx.Provide(fx.Annotate(
			repository.New,
			fx.As(new(service.Repository)),
		)),
		fx.Provide(fx.Annotate(
			service.New,
			fx.As(new(handlers.Service)),
		)),
		fx.Provide(fx.Annotate(
			handlers.New,
			fx.As(new(server.Handlers)),
		)),
		fx.Invoke(server.RegisterServer),
	).Run()
}
