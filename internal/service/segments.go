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

	seg, err := s.segmentRepo.Get(ctx, segment.Slug)
	if err != nil {
		return err
	}

	if seg != nil {
		return errors.ErrDuplicateSegment
	}

	err = s.segmentRepo.Add(ctx, segment)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SegmentDelete(ctx context.Context, slug string) error {
	if !IsValidSlug(slug) {
		return errors.ErrInvalidSegmentSlug
	}

	seg, err := s.segmentRepo.Get(ctx, slug)
	if err != nil {
		return err
	}

	if seg == nil {
		return errors.ErrSegmentNotFound
	}

	if !seg.DeletedAt.IsZero() {
		return errors.ErrAlreadyDeleted
	}

	err = s.segmentRepo.Delete(ctx, slug)
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
