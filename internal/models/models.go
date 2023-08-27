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
)
