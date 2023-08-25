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

func TestService_AddSegment(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name             string
		inputBody        models.Segment
		getSegmentReturn *models.Segment
		getSegmentError  error
		repositoryReturn error
		expectedReturn   error
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
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockRepository(ctrl)
			repo.EXPECT().GetSegment(context.Background(), tc.inputBody.Tag).Return(tc.getSegmentReturn)
			repo.EXPECT().AddSegment(context.Background(), tc.inputBody).Return(tc.repositoryReturn)

			zp, _ := zap.NewDevelopment()
			service := New(repo, zp)

			err := service.AddSegment(context.Background(), tc.inputBody)

			a.Equal(tc.expectedReturn, err)
		})
	}
}

func TestService_DeleteSegment(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name             string
		input            string
		getSegmentReturn *models.Segment
		getSegmentError  error
		repositoryReturn error
		expectedReturn   error
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
		},
		{
			name:             "Segment not found",
			input:            "NEW_TAG",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrNotFound,
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
		},
		{
			name:             "Invalid tag",
			input:            "NeW-tAg",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentTag,
		},
		{
			name:             "Empty tag",
			input:            "",
			getSegmentReturn: nil,
			getSegmentError:  nil,
			repositoryReturn: nil,
			expectedReturn:   errors.ErrInvalidSegmentTag,
		},
		{
			name:             "DB error checkDuplicate",
			input:            "NEW_TAG",
			getSegmentReturn: nil,
			getSegmentError:  pgx.ErrTxCommitRollback,
			repositoryReturn: nil,
			expectedReturn:   pgx.ErrTxCommitRollback,
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
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockRepository(ctrl)
			repo.EXPECT().GetSegment(context.Background(), tc.input).Return(tc.getSegmentReturn)
			repo.EXPECT().DeleteSegment(context.Background(), tc.input).Return(tc.repositoryReturn)

			zp, _ := zap.NewDevelopment()
			service := New(repo, zp)

			err := service.DeleteSegment(context.Background(), tc.input)

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
