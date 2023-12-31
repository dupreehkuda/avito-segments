package repository

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	errs "github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (r *Repository) SetSegments(ctx context.Context, segments *models.UserSetRequest) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	query := sq.Insert("user_segments").
		Columns("slug", "user_id", "created_at", "expired_at").
		Suffix("ON CONFLICT (slug, user_id) DO UPDATE SET expired_at = excluded.expired_at").
		PlaceholderFormat(sq.Dollar)

	for _, segment := range segments.Segments {
		query = query.Values(segment.Slug, segments.UserID, time.Now(), segment.Expire)
	}

	queryString, queryArgs := query.MustSql()
	_, err = conn.Exec(ctx, queryString, queryArgs...)
	if err != nil {
		r.logger.Error("Error while executing query", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) DeleteSegments(ctx context.Context, segments *models.UserDeleteRequest) error {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	queryString, queryArgs := sq.Update("user_segments").
		Set("deleted_at", time.Now()).
		Where(sq.Eq{
			"user_id": segments.UserID,
			"slug":    segments.Slugs,
		}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	_, err = conn.Exec(ctx, queryString, queryArgs...)
	if err != nil {
		r.logger.Error("Error while executing query", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) GetSegments(ctx context.Context, userID string) (*models.UserResponse, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	queryString, queryArgs := sq.Select("user_segments.slug").
		From("user_segments").
		Join("segments on segments.slug = user_segments.slug").
		Where(
			sq.Eq{
				"segments.deleted_at":      nil,
				"user_segments.deleted_at": nil,
				"user_segments.user_id":    userID,
			},
			sq.LtOrEq{
				"user_segments.expired_at": time.Now(),
			},
		).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	rows, err := conn.Query(ctx, queryString, queryArgs...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrSegmentNotFound
		}

		r.logger.Error("Error while executing query", zap.Error(err))
		return nil, err
	}

	resp := &models.UserResponse{
		UserID: userID,
		Slugs:  make([]string, 0),
	}

	for rows.Next() {
		var slug string

		err = rows.Scan(&slug)
		if err != nil {
			r.logger.Error("Error while scanning query", zap.Error(err))
			return nil, err
		}

		resp.Slugs = append(resp.Slugs, slug)
	}

	return resp, nil
}

func (r *Repository) GetReportData(ctx context.Context, year, month int) ([]models.ReportRow, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		r.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()

	query := `
		SELECT user_id, slug, 'added' as method, created_at AS timestamp
		FROM user_segments
		WHERE EXTRACT(YEAR FROM created_at) = $1
		  AND EXTRACT(MONTH FROM created_at) = $2
		
		UNION
		
		SELECT user_id, slug, 'deleted', deleted_at AS timestamp
		FROM user_segments
		WHERE EXTRACT(YEAR FROM deleted_at) = $1
		  AND EXTRACT(MONTH FROM deleted_at) = $2
		
		UNION
		
		SELECT user_id, slug, 'expired', expired_at AS timestamp
		FROM user_segments
		WHERE EXTRACT(YEAR FROM expired_at) = $1
		  AND EXTRACT(MONTH FROM expired_at) = $2
		
		ORDER BY timestamp;
	`

	rows, err := conn.Query(ctx, query, year, month)
	if err != nil {
		r.logger.Error("Error while executing query", zap.Error(err))
		return nil, err
	}

	resp := make([]models.ReportRow, 0)

	for rows.Next() {
		var row models.ReportRow

		err = rows.Scan(&row.UserID, &row.Slug, &row.Method, &row.Timestamp)
		if err != nil {
			r.logger.Error("Error while scanning query", zap.Error(err))
			return nil, err
		}

		resp = append(resp, row)
	}

	if len(resp) == 0 {
		return nil, errs.ErrDataNotFound
	}

	return resp, nil
}
