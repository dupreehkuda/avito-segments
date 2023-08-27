package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/dupreehkuda/avito-segments/internal/errors"
	"github.com/dupreehkuda/avito-segments/internal/models"
	"github.com/dupreehkuda/avito-segments/internal/service"
)

func TestService_UserGetSegments(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name      string
		inputBody string

		repositoryReturn *models.UserResponse
		repositoryError  error

		expectedReturn *models.UserResponse
		expectedError  error

		expectingRepositoryCall bool
	}{
		{
			name:      "Segment returned",
			inputBody: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			repositoryReturn: &models.UserResponse{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs: []string{
					"TEST_SLUG",
				},
			},
			repositoryError: nil,
			expectedReturn: &models.UserResponse{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs: []string{
					"TEST_SLUG",
				},
			},
			expectedError:           nil,
			expectingRepositoryCall: true,
		},
		{
			name:      "No segments",
			inputBody: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			repositoryReturn: &models.UserResponse{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{},
			},
			repositoryError:         nil,
			expectedReturn:          nil,
			expectedError:           errors.ErrSegmentsNotFound,
			expectingRepositoryCall: true,
		},
		{
			name:      "Nil segments",
			inputBody: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			repositoryReturn: &models.UserResponse{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  nil,
			},
			repositoryError:         nil,
			expectedReturn:          nil,
			expectedError:           errors.ErrSegmentsNotFound,
			expectingRepositoryCall: true,
		},
		{
			name:                    "Not found",
			inputBody:               "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			repositoryReturn:        nil,
			repositoryError:         errors.ErrUserNotFound,
			expectedReturn:          nil,
			expectedError:           errors.ErrUserNotFound,
			expectingRepositoryCall: true,
		},
		{
			name:                    "Some internal error",
			inputBody:               "80b0b88d-379e-11ee-8bf7-0242c0a80002",
			repositoryReturn:        nil,
			repositoryError:         os.ErrInvalid,
			expectedReturn:          nil,
			expectedError:           os.ErrInvalid,
			expectingRepositoryCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockRepository(ctrl)

			if tc.expectingRepositoryCall {
				repo.EXPECT().UserGetSegments(context.Background(), tc.inputBody).Return(tc.repositoryReturn, tc.repositoryError)
			}

			zp, _ := zap.NewDevelopment()
			serv := service.New(repo, zp)

			resp, err := serv.UserGetSegments(context.Background(), tc.inputBody)

			a.Equal(tc.expectedReturn, resp)
			a.Equal(tc.expectedError, err)
		})
	}
}

func TestService_UserSetSegments(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name      string
		inputBody *models.UserSetRequest

		getCountInput  []string
		getCountReturn int
		getCountError  error

		repositoryError error

		expectedError error

		expectingGetCountCall   bool
		expectingRepositoryCall bool
	}{
		{
			name: "Segment created",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
				},
			},
			getCountInput:           []string{"TEST_SLUG"},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           nil,
			expectingGetCountCall:   true,
			expectingRepositoryCall: true,
		},
		{
			name: "Segment created w/ expire",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug:   "TEST_SLUG",
						Expire: time.Date(2024, time.August, 26, 19, 00, 00, 00, time.Local),
					},
				},
			},
			getCountInput:           []string{"TEST_SLUG"},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           nil,
			expectingGetCountCall:   true,
			expectingRepositoryCall: true,
		},
		{
			name: "Invalid slug",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
					{
						Slug: "TEST_SLUG-2",
					},
				},
			},
			getCountReturn:          0,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           errors.ErrInvalidSegmentSlug,
			expectingGetCountCall:   false,
			expectingRepositoryCall: false,
		},
		{
			name: "Already expired",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
					{
						Slug:   "TEST_SLUG_2",
						Expire: time.Date(2022, time.August, 26, 19, 00, 00, 00, time.Local),
					},
				},
			},
			getCountReturn:          0,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           errors.ErrAlreadyExpired,
			expectingGetCountCall:   false,
			expectingRepositoryCall: false,
		},
		{
			name: "Count mismatch",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
					{
						Slug: "TEST_SLUG_2",
					},
				},
			},
			getCountInput:           []string{"TEST_SLUG", "TEST_SLUG_2"},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           errors.ErrSegmentsNotFound,
			expectingGetCountCall:   true,
			expectingRepositoryCall: false,
		},
		{
			name: "Some internal error",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
				},
			},
			getCountInput:           []string{"TEST_SLUG"},
			getCountReturn:          0,
			getCountError:           os.ErrInvalid,
			repositoryError:         nil,
			expectedError:           os.ErrInvalid,
			expectingGetCountCall:   true,
			expectingRepositoryCall: false,
		},
		{
			name: "Some internal error",
			inputBody: &models.UserSetRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Segments: []models.UserSegment{
					{
						Slug: "TEST_SLUG",
					},
				},
			},
			getCountInput:           []string{"TEST_SLUG"},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         os.ErrInvalid,
			expectedError:           os.ErrInvalid,
			expectingGetCountCall:   true,
			expectingRepositoryCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockRepository(ctrl)

			if tc.expectingGetCountCall {
				repo.EXPECT().SegmentCount(context.Background(), tc.getCountInput).Return(tc.getCountReturn, tc.getCountError)
			}

			if tc.expectingRepositoryCall {
				repo.EXPECT().UserSetSegments(context.Background(), tc.inputBody).Return(tc.repositoryError)
			}

			zp, _ := zap.NewDevelopment()
			serv := service.New(repo, zp)

			err := serv.UserSetSegments(context.Background(), tc.inputBody)

			a.Equal(tc.expectedError, err)
		})
	}
}

func TestService_UserDeleteSegments(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name      string
		inputBody *models.UserDeleteRequest

		getCountReturn int
		getCountError  error

		repositoryError error

		expectedError error

		expectingGetCountCall   bool
		expectingRepositoryCall bool
	}{
		{
			name: "Segment created",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{"TEST_SLUG"},
			},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           nil,
			expectingGetCountCall:   true,
			expectingRepositoryCall: true,
		},
		{
			name: "Invalid slug",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{"TEST_SLUG", "TEST_SLUG-2"},
			},
			getCountReturn:          0,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           errors.ErrInvalidSegmentSlug,
			expectingGetCountCall:   false,
			expectingRepositoryCall: false,
		},
		{
			name: "Count mismatch",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{"TEST_SLUG", "TEST_SLUG_2"},
			},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         nil,
			expectedError:           errors.ErrSegmentsNotFound,
			expectingGetCountCall:   true,
			expectingRepositoryCall: false,
		},
		{
			name: "Some internal error",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{"TEST_SLUG"},
			},
			getCountReturn:          0,
			getCountError:           os.ErrInvalid,
			repositoryError:         nil,
			expectedError:           os.ErrInvalid,
			expectingGetCountCall:   true,
			expectingRepositoryCall: false,
		},
		{
			name: "Some internal error",
			inputBody: &models.UserDeleteRequest{
				UserID: "80b0b88d-379e-11ee-8bf7-0242c0a80002",
				Slugs:  []string{"TEST_SLUG"},
			},
			getCountReturn:          1,
			getCountError:           nil,
			repositoryError:         os.ErrInvalid,
			expectedError:           os.ErrInvalid,
			expectingGetCountCall:   true,
			expectingRepositoryCall: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockRepository(ctrl)

			if tc.expectingGetCountCall {
				repo.EXPECT().SegmentCount(context.Background(), tc.inputBody.Slugs).Return(tc.getCountReturn, tc.getCountError)
			}

			if tc.expectingRepositoryCall {
				repo.EXPECT().UserDeleteSegments(context.Background(), tc.inputBody).Return(tc.repositoryError)
			}

			zp, _ := zap.NewDevelopment()
			serv := service.New(repo, zp)

			err := serv.UserDeleteSegments(context.Background(), tc.inputBody)

			a.Equal(tc.expectedError, err)
		})
	}
}
