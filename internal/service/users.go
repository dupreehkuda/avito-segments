package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
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

	count, err := s.segmentRepo.Count(ctx, slugs)
	if err != nil {
		return err
	}

	if len(req.Segments) != count {
		return errors.ErrSegmentsNotFound
	}

	err = s.userRepo.SetSegments(ctx, req)
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

	count, err := s.segmentRepo.Count(ctx, req.Slugs)
	if err != nil {
		return err
	}

	if len(req.Slugs) != count {
		return errors.ErrSegmentsNotFound
	}

	err = s.userRepo.DeleteSegments(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UserGetSegments(ctx context.Context, userID string) (*models.UserResponse, error) {
	resp, err := s.userRepo.GetSegments(ctx, userID)
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

func (s *Service) CreateReport(ctx context.Context, year, month int) (string, error) {
	if !IsValidReportTime(year, month) {
		return "", errors.ErrInvalidPeriod
	}

	data, err := s.userRepo.GetReportData(ctx, year, month)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%v_%v_report.csv", month, year)

	file, err := os.Create("reports/" + fileName)
	if err != nil {
		return "", fmt.Errorf("error creating CSV: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		err = writer.Write([]string{row.UserID, row.Slug, row.Method, row.Timestamp.String()})
		if err != nil {
			return "", fmt.Errorf("error writing to CSV: %w", err)
		}
	}

	return fileName, nil
}

func IsValidReportTime(year, month int) bool {
	now := time.Now()

	if year <= 1971 || year >= now.Year() {
		return false
	}

	if month < 1 || month > 12 {
		return false
	}

	return true
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
