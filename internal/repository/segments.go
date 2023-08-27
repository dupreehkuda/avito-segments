package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (r *Repository) Add(ctx context.Context, segment *models.Segment) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	queryString, queryArgs := sq.Insert("segments").Columns("slug", "description", "created_at").
		Values(segment.Slug, segment.Description, time.Now()).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	if _, err = conn.Exec(ctx, queryString, queryArgs...); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, slug string) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	queryString, queryArgs := sq.Update("segments").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"slug": slug}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	if _, err = conn.Exec(ctx, queryString, queryArgs...); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, slug string) (*models.Segment, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	res := &models.Segment{}
	var deletedAt sql.NullTime

	queryString, queryArgs := sq.Select("slug", "description", "deleted_at").
		From("segments").
		Where(sq.Eq{"slug": slug}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	err = conn.QueryRow(ctx, queryString, queryArgs...).
		Scan(&res.Slug, &res.Description, &deletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	res.DeletedAt = deletedAt.Time

	return res, nil
}

func (r *Repository) Count(ctx context.Context, slugs []string) (int, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return 0, err
	}
	defer conn.Release()

	var count int

	queryString, queryArgs := sq.Select("COUNT(*)").
		From("segments").
		Where(sq.Eq{"slug": slugs}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	err = conn.QueryRow(ctx, queryString, queryArgs...).
		Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
