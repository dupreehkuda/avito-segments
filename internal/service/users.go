package service

import (
	"context"
	"time"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (s *Service) UserSetSegments(ctx context.Context, req *models.UserSetRequest) error {
	var slugs = make([]string, 0, len(req.Segments))

	for _, segment := range req.Segments {
		if err := IsValidSegment(segment); err != nil {
			return err
		}

		slugs = append(slugs, segment.Slug)
	}

	count, err := s.repository.SegmentCount(ctx, slugs)
	if err != nil {
		return err
	}

	if len(req.Segments) != count {
		return errors.ErrSegmentsNotFound
	}

	err = s.repository.UserSetSegments(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UserDeleteSegments(ctx context.Context, req *models.UserDeleteRequest) error {
	for _, segment := range req.Slugs {
		if !IsValidSlug(segment) {
			return errors.ErrInvalidSegmentSlug
		}
	}

	count, err := s.repository.SegmentCount(ctx, req.Slugs)
	if err != nil {
		return err
	}

	if len(req.Slugs) != count {
		return errors.ErrSegmentsNotFound
	}

	err = s.repository.UserDeleteSegments(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UserGetSegments(ctx context.Context, userID string) (*models.UserResponse, error) {
	resp, err := s.repository.UserGetSegments(ctx, userID)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, errors.ErrUserNotFound
	}

	if resp.Slugs == nil || len(resp.Slugs) == 0 {
		return nil, errors.ErrSegmentsNotFound
	}

	return resp, nil
}

func IsValidSegment(segment models.UserSegment) error {
	if !IsValidSlug(segment.Slug) {
		return errors.ErrInvalidSegmentSlug
	}

	if segment.Expire.IsZero() {
		return nil
	}

	if segment.Expire.Before(time.Now()) {
		return errors.ErrAlreadyExpired
	}

	return nil
}
