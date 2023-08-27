package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
	"github.com/dupreehkuda/avito-segments/internal/service"
)

func TestService_SegmentAdd(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name             string
		inputBody        *models.Segment
		getSegmentReturn *models.Segment
		getSegmentError  error
		repositoryReturn error
		expectedReturn   error
		expectingGet     bool
		expectingAdd     bool
	}{
		{
			name: "Segment created",
			inputBody: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   nil,
			expectingGet:     true,
			expectingAdd:     true,
		},
		{
			name: "Invalid slug",
			inputBody: &models.Segment{
				Slug: "NeW-slUg",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentSlug,
			expectingGet:     false,
			expectingAdd:     false,
		},
		{
			name: "Empty slug",
			inputBody: &models.Segment{
				Slug: "",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentSlug,
			expectingGet:     false,
			expectingAdd:     false,
		},
		{
			name: "Duplicate entry",
			inputBody: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
			},
			getSegmentReturn: &models.Segment{
				Slug:        "OLD_SLUG",
				Description: "just old slug",
				DeletedAt:   time.Time{},
			},
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrDuplicateSegment,
			expectingGet:     true,
			expectingAdd:     false,
		},
		{
			name: "DB error checkDuplicate",
			inputBody: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
			},
			getSegmentReturn: nil,
			getSegmentError:  pgx.ErrTxCommitRollback,
			repositoryReturn: nil,
			expectedReturn:   pgx.ErrTxCommitRollback,
			expectingGet:     true,
			expectingAdd:     false,
		},
		{
			name: "DB error addSegment",
			inputBody: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: pgx.ErrTxClosed,
			expectedReturn:   pgx.ErrTxClosed,
			expectingGet:     true,
			expectingAdd:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := NewMockUserRepository(ctrl)
			segmentRepo := NewMockSegmentRepository(ctrl)

			if tc.expectingGet {
				segmentRepo.EXPECT().Get(context.Background(), tc.inputBody.Slug).Return(tc.getSegmentReturn, tc.getSegmentError)
			}

			if tc.expectingAdd {
				segmentRepo.EXPECT().Add(context.Background(), tc.inputBody).Return(tc.repositoryReturn)
			}

			zp, _ := zap.NewDevelopment()
			serv := service.New(userRepo, segmentRepo, zp)

			err := serv.SegmentAdd(context.Background(), tc.inputBody)

			a.Equal(tc.expectedReturn, err)
		})
	}
}

func TestService_SegmentDelete(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name             string
		input            string
		getSegmentReturn *models.Segment
		getSegmentError  error
		repositoryReturn error
		expectedReturn   error
		expectingGet     bool
		expectingDelete  bool
	}{
		{
			name:  "Segment deleted",
			input: "NEW_SLUG",
			getSegmentReturn: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
				DeletedAt:   time.Time{},
			},
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   nil,
			expectingGet:     true,
			expectingDelete:  true,
		},
		{
			name:             "Segment not found",
			input:            "NEW_SLUG",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrSegmentNotFound,
			expectingGet:     true,
			expectingDelete:  false,
		},
		{
			name:  "Segment already deleted",
			input: "NEW_SLUG",
			getSegmentReturn: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
				DeletedAt:   time.Date(2022, time.August, 26, 19, 00, 00, 00, time.Local),
			},
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrAlreadyDeleted,
			expectingGet:     true,
			expectingDelete:  false,
		},
		{
			name:             "Invalid slug",
			input:            "NeW-SLug",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentSlug,
			expectingGet:     false,
			expectingDelete:  false,
		},
		{
			name:             "Empty slug",
			input:            "",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentSlug,
			expectingGet:     false,
			expectingDelete:  false,
		},
		{
			name:             "DB error checkDuplicate",
			input:            "NEW_SLUG",
			getSegmentReturn: nil,
			getSegmentError:  pgx.ErrTxCommitRollback,
			repositoryReturn: nil,
			expectedReturn:   pgx.ErrTxCommitRollback,
			expectingGet:     true,
			expectingDelete:  false,
		},
		{
			name:  "DB error addSegment",
			input: "NEW_SLUG",
			getSegmentReturn: &models.Segment{
				Slug:        "NEW_SLUG",
				Description: "just new slug",
				DeletedAt:   time.Time{},
			},
			getSegmentError:  nil,
			repositoryReturn: pgx.ErrTxClosed,
			expectedReturn:   pgx.ErrTxClosed,
			expectingGet:     true,
			expectingDelete:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := NewMockUserRepository(ctrl)
			segmentRepo := NewMockSegmentRepository(ctrl)

			if tc.expectingGet {
				segmentRepo.EXPECT().Get(context.Background(), tc.input).Return(tc.getSegmentReturn, tc.getSegmentError)
			}

			if tc.expectingDelete {
				segmentRepo.EXPECT().Delete(context.Background(), tc.input).Return(tc.repositoryReturn)
			}

			zp, _ := zap.NewDevelopment()
			serv := service.New(userRepo, segmentRepo, zp)

			err := serv.SegmentDelete(context.Background(), tc.input)

			a.Equal(tc.expectedReturn, err)
		})
	}
}

func TestIsValidSlug(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name string
		slug string
		want bool
	}{
		{
			name: "Valid string w/out numbers",
			slug: "AVITO_PERFORMANCE_VAS",
			want: true,
		},
		{
			name: "Valid string w/ numbers",
			slug: "AVITO_DISCOUNT_50",
			want: true,
		},
		{
			name: "Invalid w/ dashes",
			slug: "AVITO-PERFORMANCE-VAS",
			want: false,
		},
		{
			name: "Invalid w/ other chars",
			slug: "reAl-tEst$thO_",
			want: false,
		},
		{
			name: "Empty string",
			slug: "",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a.Equal(tc.want, service.IsValidSlug(tc.slug))
		})
	}
}

func TestIsValidSegment(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name    string
		segment models.UserSegment
		want    error
	}{
		{
			name: "Valid date",
			segment: models.UserSegment{
				Slug:   "AVITO_PERFORMANCE_VAS",
				Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.Local),
			},
			want: nil,
		},
		{
			name: "Past date",
			segment: models.UserSegment{
				Slug:   "AVITO_PERFORMANCE_VAS",
				Expire: time.Date(2022, time.August, 26, 19, 00, 00, 00, time.Local),
			},
			want: errors.ErrAlreadyExpired,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a.Equal(tc.want, service.IsValidSegment(tc.segment))
		})
	}
}
