package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/config"
)

// Repository provides a connection with database
type Repository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// New creates a new instance of the database layer and migrates it
func New(config *config.Config, logger *zap.Logger) *Repository {
	uri := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Settings,
	)

	dbConfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		logger.Error("Unable to parse config", zap.Error(err))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Error("Unable to connect to database", zap.Error(err))
	}

	schema, err := os.ReadFile(config.Database.MigrationPath)
	if err != nil {
		logger.Error("Error occurred while getting migration schema", zap.Error(err))
	}

	_, err = pool.Exec(context.Background(), string(schema))
	if err != nil {
		logger.Error("Error occurred while executing schema", zap.Error(err))
	}

	return &Repository{
		pool:   pool,
		logger: logger,
	}
}