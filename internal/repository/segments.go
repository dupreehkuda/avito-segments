package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
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

	queryString, queryArgs := sq.Insert("segments").Columns("tag", "description", "created_at").
		Values(segment.Tag, segment.Description, segment.DeletedAt).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	if _, err = conn.Exec(ctx, queryString, queryArgs...); err != nil {
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

	queryString, queryArgs := sq.Update("segments").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"tag": tag}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	if _, err = conn.Exec(ctx, queryString, queryArgs...); err != nil {
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

	queryString, queryArgs := sq.Select("tag", "description", "deleted_at").
		From("segments").
		Where(sq.Eq{"tag": tag}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	err = conn.QueryRow(ctx, queryString, queryArgs...).
		Scan(&res.Tag, &res.Description, &res.DeletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}
