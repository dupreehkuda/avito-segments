package logger

import (
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
)

// New initializes new Logger instance.
func New(config *config.Config) *zap.Logger {
	var logger *zap.Logger

	switch config.Common.Logger {
	case "debug":
		logger, _ = zap.NewDevelopment()
	default:
		logger, _ = zap.NewProduction()
	}

	return logger
}
