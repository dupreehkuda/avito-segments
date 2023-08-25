//go:generate easyjson -no_std_marshalers models.go
package models

//easyjson:json
type (
	Segment struct {
		Tag         string `json:"tag"`
		Description string `json:"description"`
	}

	UserSegment struct {
		Tag    string `json:"tag"`
		Expire string `json:"expire"`
	}

	UserRequest struct {
		ID       string        `json:"id"`
		Segments []UserSegment `json:"segments"`
	}

	UserResponse struct {
		ID       string   `json:"id"`
		Segments []string `json:"segments"`
	}

	Error struct {
		Message string `json:"message"`
	}
)
