package errors

import "errors"

var (
	ErrDuplicateSegment  = errors.New("segment already exists")
	ErrInvalidSegmentTag = errors.New("invalid segment tag naming")
	ErrNotFound          = errors.New("segment not found")
	ErrAlreadyDeleted    = errors.New("segment had been already deleted")
)
