package service

import (
	"context"
	"regexp"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (s *Service) SegmentAdd(ctx context.Context, segment models.Segment) error {
	if !IsValidTag(segment.Tag) {
		return errors.ErrInvalidSegmentTag
	}

	seg, err := s.repository.SegmentGet(ctx, segment.Tag)
	if err != nil {
		return err
	}

	if seg != nil {
		return errors.ErrDuplicateSegment
	}

	if err = s.repository.SegmentAdd(ctx, segment); err != nil {
		return err
	}

	return nil
}

func (s *Service) SegmentDelete(ctx context.Context, tag string) error {
	if !IsValidTag(tag) {
		return errors.ErrInvalidSegmentTag
	}

	seg, err := s.repository.SegmentGet(ctx, tag)
	if err != nil {
		return err
	}

	if seg == nil {
		return errors.ErrNotFound
	}

	if seg.DeletedAt != nil {
		return errors.ErrAlreadyDeleted
	}

	if err = s.repository.SegmentDelete(ctx, tag); err != nil {
		return err
	}

	return nil
}

func IsValidTag(tag string) bool {
	pattern := "^[A-Z0-9_]+$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(tag)
}
