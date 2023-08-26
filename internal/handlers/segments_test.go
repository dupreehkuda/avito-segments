package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
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
		inputBody          models.Segment
		serviceReturn      error
		expectedStatusCode int
	}{
		{
			name: "Segment created",
			inputBody: models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
			},
			serviceReturn:      nil,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Duplicate tag",
			inputBody: models.Segment{
				Tag:         "NEW_TAG",
				Description: "duplicate tag",
			},
			serviceReturn:      errors.ErrDuplicateSegment,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid tag naming",
			inputBody: models.Segment{
				Tag: "NeW_TaG-1",
			},
			serviceReturn:      errors.ErrInvalidSegmentTag,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Empty tag",
			inputBody: models.Segment{
				Tag: "",
			},
			serviceReturn:      errors.ErrInvalidSegmentTag,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Some internal error",
			inputBody: models.Segment{
				Tag: "NeW_TaG-1",
			},
			serviceReturn:      os.ErrInvalid,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := json.Marshal(tc.inputBody)

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
			input:              "NEW_TAG",
			serviceReturn:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid tag naming",
			input:              "NeW_TaG-1",
			serviceReturn:      errors.ErrInvalidSegmentTag,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Empty tag",
			input:              "",
			serviceReturn:      errors.ErrInvalidSegmentTag,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Segment not found",
			input:              "OLD_TAG",
			serviceReturn:      errors.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Segment has been deleted",
			input:              "NEW_TAG",
			serviceReturn:      errors.ErrAlreadyDeleted,
			expectedStatusCode: http.StatusGone,
		},
		{
			name:               "Some internal error",
			input:              "NeW_TaG-1",
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
			c.SetPath("/api/v1/segment/:tag")
			c.SetParamNames("tag")
			c.SetParamValues(tc.input)

			err := server.SegmentDelete(c)
			e.DefaultHTTPErrorHandler(err, c)

			a.Equal(tc.expectedStatusCode, rec.Code, "Wrong status code")
		})
	}
}
