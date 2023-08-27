//go:generate easyjson -no_std_marshalers models.go
package models

import "time"

//easyjson:json
type (
	Segment struct {
		Slug        string `json:"slug"`
		Description string `json:"description,omitempty"`
		DeletedAt   time.Time
	}

	UserSegment struct {
		Slug   string    `json:"slug"`
		Expire time.Time `json:"expire,omitempty"`
	}

	UserSetRequest struct {
		UserID   string        `json:"userID"`
		Segments []UserSegment `json:"segments"`
	}

	UserDeleteRequest struct {
		UserID string   `json:"userID"`
		Slugs  []string `json:"slugs"`
	}

	UserResponse struct {
		UserID string   `json:"userID"`
		Slugs  []string `json:"slugs"`
	}

	ReportRow struct {
		UserID    string    `json:"userID"`
		Slug      string    `json:"slug"`
		Method    string    `json:"method"`
		Timestamp time.Time `json:"timestamp"`
	}

	ReportRequest struct {
		Year  int `json:"year"`
		Month int `json:"month"`
	}

	ReportResponse struct {
		Link string `json:"link"`
	}
)
