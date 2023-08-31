package errors

import "errors"

var (
	ErrDuplicateSegment   = errors.New("segment already exists")
	ErrInvalidSegmentSlug = errors.New("invalid segment slug naming")
	ErrSegmentNotFound    = errors.New("segment not found")
	ErrAlreadyDeleted     = errors.New("segment had been already deleted")
	ErrNoSegmentsProvided = errors.New("no segments provided in request")

	ErrInvalidUserID    = errors.New("invalid user id")
	ErrUserNotFound     = errors.New("user not found")
	ErrSegmentsNotFound = errors.New("segment(s) not found")
	ErrAlreadyExpired   = errors.New("provided segment expired")

	ErrDataNotFound   = errors.New("no data found")
	ErrInvalidPeriod  = errors.New("provided invalid period")
	ErrReportNotFound = errors.New("requested report not found")
)
