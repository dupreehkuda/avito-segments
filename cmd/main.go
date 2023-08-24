package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
	"github.com/dupreehkuda/avito-segments/internal/logger"
	"github.com/dupreehkuda/avito-segments/internal/repository"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(config.New),
		fx.Provide(logger.New),
		fx.Provide(repository.New),
		fx.Invoke(func(*repository.Repository) {}),
	).Run()
}
