package errors

import "errors"

var (
	ErrDuplicateSegment   = errors.New("segment already exists")
	ErrInvalidSegmentSlug = errors.New("invalid segment slug naming")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrSegmentNotFound    = errors.New("segment not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrAlreadyDeleted     = errors.New("segment had been already deleted")
	ErrAlreadyExpired     = errors.New("provided segment expired")
	ErrSegmentsNotFound   = errors.New("segment(s) not found")
)
