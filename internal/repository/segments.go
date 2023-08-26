package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (r *Repository) SegmentAdd(ctx context.Context, segment models.Segment) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	if _, err = conn.Exec(
		ctx,
		"INSERT INTO segments (tag, description, created_at) VALUES ($1, $2, $3)",
		segment.Tag,
		segment.Description,
		segment.DeletedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *Repository) SegmentDelete(ctx context.Context, tag string) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	if _, err = conn.Exec(ctx, "UPDATE segments SET deleted_at = now() WHERE tag = $1", tag); err != nil {
		return err
	}

	return nil
}

func (r *Repository) SegmentGet(ctx context.Context, tag string) (*models.Segment, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	res := &models.Segment{}

	err = conn.QueryRow(ctx, "SELECT tag, description, deleted_at FROM segments WHERE tag = $1;", tag).
		Scan(&res.Tag, &res.Description, &res.DeletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}
