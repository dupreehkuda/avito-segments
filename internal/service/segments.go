package service

import (
	"context"
	"regexp"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (s *Service) SegmentAdd(ctx context.Context, segment *models.Segment) error {
	if !IsValidSlug(segment.Slug) {
		return errors.ErrInvalidSegmentSlug
	}

	seg, err := s.repository.SegmentGet(ctx, segment.Slug)
	if err != nil {
		return err
	}

	if seg != nil {
		return errors.ErrDuplicateSegment
	}

	err = s.repository.SegmentAdd(ctx, segment)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SegmentDelete(ctx context.Context, slug string) error {
	if !IsValidSlug(slug) {
		return errors.ErrInvalidSegmentSlug
	}

	seg, err := s.repository.SegmentGet(ctx, slug)
	if err != nil {
		return err
	}

	if seg == nil {
		return errors.ErrSegmentNotFound
	}

	if !seg.DeletedAt.IsZero() {
		return errors.ErrAlreadyDeleted
	}

	err = s.repository.SegmentDelete(ctx, slug)
	if err != nil {
		return err
	}

	return nil
}

func IsValidSlug(slug string) bool {
	pattern := "^[A-Z0-9_]+$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(slug)
}
