package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/handlers"
	"github.com/dupreehkuda/avito-segments/internal/models"
)

func TestHandlers_SegmentAdd(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		inputBody          *models.Segment
		serviceReturn      error
		expectedStatusCode int
	}{
		{
			name: "Segment created",
			inputBody: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
			},
			serviceReturn:      nil,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Duplicate slug",
			inputBody: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "duplicate slug",
			},
			serviceReturn:      errors.ErrDuplicateSegment,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid slug naming",
			inputBody: &models.Segment{
				Slug: "NeW_SluG-1",
			},
			serviceReturn:      errors.ErrInvalidSegmentSlug,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Empty slug",
			inputBody: &models.Segment{
				Slug: "",
			},
			serviceReturn:      errors.ErrInvalidSegmentSlug,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Some internal error",
			inputBody: &models.Segment{
				Slug: "NeW_SLug-1",
			},
			serviceReturn:      os.ErrInvalid,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := easyjson.Marshal(tc.inputBody)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			service.EXPECT().SegmentAdd(context.Background(), tc.inputBody).Return(tc.serviceReturn)

			zp, _ := zap.NewDevelopment()
			server := handlers.New(service, zp)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/segment")

			err := server.SegmentAdd(c)
			e.DefaultHTTPErrorHandler(err, c)

			a.Equal(tc.expectedStatusCode, rec.Code, "Wrong status code")
		})
	}
}

func TestHandlers_SegmentDelete(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		input              string
		serviceReturn      error
		expectedStatusCode int
	}{
		{
			name:               "Segment created",
			input:              "NEW_SLUG",
			serviceReturn:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid slug naming",
			input:              "NeW_SLug-1",
			serviceReturn:      errors.ErrInvalidSegmentSlug,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Empty slug",
			input:              "",
			serviceReturn:      errors.ErrInvalidSegmentSlug,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Segment not found",
			input:              "OLD_SLUG",
			serviceReturn:      errors.ErrSegmentNotFound,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Segment has been deleted",
			input:              "NEW_SLUG",
			serviceReturn:      errors.ErrAlreadyDeleted,
			expectedStatusCode: http.StatusGone,
		},
		{
			name:               "Some internal error",
			input:              "NeW_SLug-1",
			serviceReturn:      os.ErrInvalid,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := NewMockService(ctrl)
			service.EXPECT().SegmentDelete(context.Background(), tc.input).Return(tc.serviceReturn)

			zp, _ := zap.NewDevelopment()
			server := handlers.New(service, zp)

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/segment/:slug")
			c.SetParamNames("slug")
			c.SetParamValues(tc.input)

			err := server.SegmentDelete(c)
			e.DefaultHTTPErrorHandler(err, c)

			a.Equal(tc.expectedStatusCode, rec.Code, "Wrong status code")
		})
	}
}
