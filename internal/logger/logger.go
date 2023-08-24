package logger

import (
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
)

var Logger *zap.Logger

// New initializes new Logger instance
func New(config *config.Config) *zap.Logger {
	switch config.Common.Logger {
	case "debug":
		Logger, _ = zap.NewDevelopment()
	default:
		Logger, _ = zap.NewProduction()
	}

	return Logger
}
