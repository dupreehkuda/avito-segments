package service

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
)

func TestService_SegmentAdd(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name             string
		inputBody        models.Segment
		getSegmentReturn *models.Segment
		getSegmentError  error
		repositoryReturn error
		expectedReturn   error
		expectingGet     bool
		expectingAdd     bool
	}{
		{
			name: "Segment created",
			inputBody: models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   nil,
			expectingGet:     true,
			expectingAdd:     true,
		},
		{
			name: "Invalid tag",
			inputBody: models.Segment{
				Tag: "NeW-tAg",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentTag,
			expectingGet:     false,
			expectingAdd:     false,
		},
		{
			name: "Empty tag",
			inputBody: models.Segment{
				Tag: "",
			},
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentTag,
			expectingGet:     false,
			expectingAdd:     false,
		},
		{
			name: "Duplicate entry",
			inputBody: models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
			},
			getSegmentReturn: &models.Segment{
				Tag:         "OLD_TAG",
				Description: "just old tag",
				DeletedAt:   nil,
			},
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrDuplicateSegment,
			expectingGet:     true,
			expectingAdd:     false,
		},
		{
			name: "DB error checkDuplicate",
			inputBody: models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
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
			inputBody: models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
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

			repo := NewMockRepository(ctrl)

			if tc.expectingGet {
				repo.EXPECT().SegmentGet(context.Background(), tc.inputBody.Tag).Return(tc.getSegmentReturn, tc.getSegmentError)
			}

			if tc.expectingAdd {
				repo.EXPECT().SegmentAdd(context.Background(), tc.inputBody).Return(tc.repositoryReturn)
			}

			zp, _ := zap.NewDevelopment()
			service := New(repo, zp)

			err := service.SegmentAdd(context.Background(), tc.inputBody)

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
			input: "NEW_TAG",
			getSegmentReturn: &models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
				DeletedAt:   nil,
			},
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   nil,
			expectingGet:     true,
			expectingDelete:  true,
		},
		{
			name:             "Segment not found",
			input:            "NEW_TAG",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrNotFound,
			expectingGet:     true,
			expectingDelete:  false,
		},
		{
			name:  "Segment already deleted",
			input: "NEW_TAG",
			getSegmentReturn: &models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
				DeletedAt:   &time.Time{},
			},
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrAlreadyDeleted,
			expectingGet:     true,
			expectingDelete:  false,
		},
		{
			name:             "Invalid tag",
			input:            "NeW-tAg",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentTag,
			expectingGet:     false,
			expectingDelete:  false,
		},
		{
			name:             "Empty tag",
			input:            "",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentTag,
			expectingGet:     false,
			expectingDelete:  false,
		},
		{
			name:             "DB error checkDuplicate",
			input:            "NEW_TAG",
			getSegmentReturn: nil,
			getSegmentError:  pgx.ErrTxCommitRollback,
			repositoryReturn: nil,
			expectedReturn:   pgx.ErrTxCommitRollback,
			expectingGet:     true,
			expectingDelete:  false,
		},
		{
			name:  "DB error addSegment",
			input: "NEW_TAG",
			getSegmentReturn: &models.Segment{
				Tag:         "NEW_TAG",
				Description: "just new tag",
				DeletedAt:   nil,
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

			repo := NewMockRepository(ctrl)

			if tc.expectingGet {
				repo.EXPECT().SegmentGet(context.Background(), tc.input).Return(tc.getSegmentReturn, tc.getSegmentError)
			}

			if tc.expectingDelete {
				repo.EXPECT().SegmentDelete(context.Background(), tc.input).Return(tc.repositoryReturn)
			}

			zp, _ := zap.NewDevelopment()
			service := New(repo, zp)

			err := service.SegmentDelete(context.Background(), tc.input)

			a.Equal(tc.expectedReturn, err)
		})
	}
}

func Test_isValidTag(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		name string
		tag  string
		want bool
	}{
		{
			name: "Valid string w/out numbers",
			tag:  "AVITO_PERFORMANCE_VAS",
			want: true,
		},
		{
			name: "Valid string w/ numbers",
			tag:  "AVITO_DISCOUNT_50",
			want: true,
		},
		{
			name: "Invalid w/ dashes",
			tag:  "AVITO-PERFORMANCE-VAS",
			want: false,
		},
		{
			name: "Invalid w/ other chars",
			tag:  "reAl-tEst$thO_",
			want: false,
		},
		{
			name: "Empty string",
			tag:  "",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			a.Equal(tc.want, isValidTag(tc.tag))
		})
	}
}
