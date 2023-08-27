package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/handlers"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func TestHandlers_UserSetSegments(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name                 string
		inputBody            *models.UserSetRequest
		expectingServiceCall bool
		serviceReturn        error
		expectedStatusCode   int
	}{
		{
			name: "Segment added",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name: "Segment added w/ expiration",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name: "Segment already expired",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2022, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        errors.ErrAlreadyExpired,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Invalid slug",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "tEsT_sLUg",
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        errors.ErrInvalidSegmentSlug,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "One of invalid slug",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
					{
						Slug: "tEsT_sLUg",
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        errors.ErrInvalidSegmentSlug,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "One of already expired",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
					{
						Slug:   "TEST_SLUG_2",
						Expire: time.Date(2022, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        errors.ErrAlreadyExpired,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Some internal error",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
				},
			},
			expectingServiceCall: true,
			serviceReturn:        os.ErrInvalid,
			expectedStatusCode:   http.StatusInternalServerError,
		},
		{
			name: "Invalid user id",
			inputBody: &models.UserSetRequest{
				UserID: "123456",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.UTC),
					},
				},
			},
			expectingServiceCall: false,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "No segments sent",
			inputBody: &models.UserSetRequest{
				UserID:   "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: nil,
			},
			expectingServiceCall: false,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := easyjson.Marshal(tc.inputBody)
			var input models.UserSetRequest
			_ = easyjson.Unmarshal(data, &input)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			if tc.expectingServiceCall {
				service.EXPECT().UserSetSegments(context.Background(), &input).Return(tc.serviceReturn)
			}

			zp, _ := zap.NewDevelopment()
			server := handlers.New(service, zp)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/user")

			err := server.UserSetSegments(c)
			e.DefaultHTTPErrorHandler(err, c)

			a.Equal(tc.expectedStatusCode, rec.Code, "Wrong status code")
		})
	}
}

func TestHandlers_UserDeleteSegments(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name                 string
		inputBody            *models.UserDeleteRequest
		expectingServiceCall bool
		serviceReturn        error
		expectedStatusCode   int
	}{
		{
			name: "Segment deleted",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs: []string{
					"TEST_SLUG",
				},
			},
			expectingServiceCall: true,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name: "Segments deleted",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs: []string{
					"TEST_SLUG",
					"TEST_SLUG_2",
				},
			},
			expectingServiceCall: true,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusOK,
		},
		{
			name: "Invalid user id",
			inputBody: &models.UserDeleteRequest{
				UserID: "123456",
				Slugs: []string{
					"TEST_SLUG",
				},
			},
			expectingServiceCall: false,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "No slugs provided",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  nil,
			},
			expectingServiceCall: false,
			serviceReturn:        nil,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Invalid slug",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs: []string{
					"TEST_SLUG",
					"TeSt_sLuG-1",
				},
			},
			expectingServiceCall: true,
			serviceReturn:        errors.ErrInvalidSegmentSlug,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name: "Some internal service",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs: []string{
					"TEST_SLUG",
				},
			},
			expectingServiceCall: true,
			serviceReturn:        os.ErrInvalid,
			expectedStatusCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := easyjson.Marshal(tc.inputBody)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			if tc.expectingServiceCall {
				service.EXPECT().UserDeleteSegments(context.Background(), tc.inputBody).Return(tc.serviceReturn)
			}

			zp, _ := zap.NewDevelopment()
			server := handlers.New(service, zp)

			req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewReader(data))
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/user")

			err := server.UserDeleteSegments(c)
			e.DefaultHTTPErrorHandler(err, c)

			a.Equal(tc.expectedStatusCode, rec.Code, "Wrong status code")
		})
	}
}

func TestHandlers_UserGetSegments(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name                 string
		input                string
		expectingServiceCall bool
		serviceReturnData    *models.UserResponse
		serviceReturnError   error
		expectedStatusCode   int
	}{
		{
			name:                 "Successful request",
			input:                "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			expectingServiceCall: true,
			serviceReturnData: &models.UserResponse{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{"TEST_SLUG"},
			},
			serviceReturnError: nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:                 "Invalid user id",
			input:                "123456",
			expectingServiceCall: false,
			serviceReturnData:    nil,
			serviceReturnError:   nil,
			expectedStatusCode:   http.StatusBadRequest,
		},
		{
			name:                 "No segments",
			input:                "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			expectingServiceCall: true,
			serviceReturnData:    nil,
			serviceReturnError:   errors.ErrSegmentsNotFound,
			expectedStatusCode:   http.StatusNoContent,
		},
		{
			name:                 "Invalid not found",
			input:                "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			expectingServiceCall: true,
			serviceReturnData:    nil,
			serviceReturnError:   errors.ErrUserNotFound,
			expectedStatusCode:   http.StatusNotFound,
		},
		{
			name:                 "Some internal error",
			input:                "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			expectingServiceCall: true,
			serviceReturnData:    nil,
			serviceReturnError:   os.ErrInvalid,
			expectedStatusCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			if tc.expectingServiceCall {
				service.EXPECT().UserGetSegments(context.Background(), tc.input).Return(tc.serviceReturnData, tc.serviceReturnError)
			}

			zp, _ := zap.NewDevelopment()
			server := handlers.New(service, zp)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/user/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.input)

			err := server.UserGetSegments(c)
			e.DefaultHTTPErrorHandler(err, c)

			a.Equal(tc.expectedStatusCode, rec.Code, "Wrong status code")
		})
	}
}
