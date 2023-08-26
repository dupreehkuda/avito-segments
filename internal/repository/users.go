package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (r *Repository) UserSetSegments(ctx context.Context, segments models.UserRequest) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	query := sq.Insert("user_segments").
		Columns("tag", "user_id", "created_at", "expired_at").
		Suffix("ON CONFLICT DO UPDATE SET expired_at = excluded.expired_at").
		PlaceholderFormat(sq.Dollar)

	for _, segment := range segments.Segments {
		query.Values(segment.Tag, segments.ID, time.Now(), segment.Expire)
	}

	queryString, queryArgs := query.MustSql()
	_, err = conn.Exec(ctx, queryString, queryArgs...)
	if err != nil {
		r.logger.Error("Error while executing query", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) UserDeleteSegments(ctx context.Context, userID string, segments []string) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	queryString, queryArgs := sq.Update("user_segments").
		Set("deleted_at", time.Now()).
		Where(sq.Eq{
			"user_id": userID,
			"tag":     segments,
		}).
		MustSql()

	_, err = conn.Exec(ctx, queryString, queryArgs...)
	if err != nil {
		r.logger.Error("Error while executing query", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) UserGetSegments(ctx context.Context, userID string) (*models.UserResponse, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	queryString, queryArgs := sq.Select("user_segments.tag").
		From("user_segments").
		Join("segments on segments.tag = user_segments.tag").
		Where(sq.Eq{
			"segments.deleted_at":      nil,
			"user_segments.deleted_at": nil,
			"user_segments.user_id":    userID,
		}).
		Suffix("AND user_segments.expired_at > now() OR user_segments.expired_at IS NULL").
		MustSql()

	rows, err := conn.Query(ctx, queryString, queryArgs...)
	if err != nil {
		r.logger.Error("Error while executing query", zap.Error(err))
		return nil, err
	}

	resp := &models.UserResponse{
		ID:       userID,
		Segments: make([]string, 0),
	}

	for rows.Next() {
		var slug string

		err = rows.Scan(&slug)
		if err != nil {
			r.logger.Error("Error while scanning query", zap.Error(err))
			return nil, err
		}

		resp.Segments = append(resp.Segments, slug)
	}

	return resp, nil
}
