package service

import (
	"context"
	"regexp"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func (s *Service) AddSegment(ctx context.Context, segment models.Segment) error {
	if !isValidTag(segment.Tag) {
		return errors.ErrInvalidSegmentTag
	}

	return nil
}

func (s *Service) DeleteSegment(ctx context.Context, tag string) error {
	if !isValidTag(tag) {
		return errors.ErrInvalidSegmentTag
	}

	return nil
}

func isValidTag(tag string) bool {
	pattern := "^[A-Z0-9_]+$"
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(tag)
}
