package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zhwjimmy/user-center/internal/dto"
	"github.com/zhwjimmy/user-center/internal/mock"
	"github.com/zhwjimmy/user-center/internal/model"
	"go.uber.org/zap"
)

func strPtr(s string) *string { return &s }

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		user          *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name: "successful user creation",
			user: &model.User{
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				FirstName:    strPtr("Test"),
				LastName:     strPtr("User"),
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				// Check if user with email already exists
				repo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, assert.AnError) // User not found by email
				// Check if user with username already exists
				repo.EXPECT().GetByUsername(gomock.Any(), "testuser").
					Return(nil, assert.AnError) // User not found by username
				// Create user
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&model.User{
					ID:           "test-user-id",
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: "hashedpassword",
					FirstName:    strPtr("Test"),
					LastName:     strPtr("User"),
					IsActive:     true,
				}, nil)
			},
		},
		{
			name: "database error",
			user: &model.User{
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
			},
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				// Check if user with email already exists
				repo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
					Return(nil, assert.AnError) // User not found by email
				// Check if user with username already exists
				repo.EXPECT().GetByUsername(gomock.Any(), "testuser").
					Return(nil, assert.AnError) // User not found by username
				// Create user fails
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			result, err := service.CreateUser(context.Background(), tt.user)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.user.Username, result.Username)
				assert.Equal(t, tt.user.Email, result.Email)
				assert.True(t, result.IsActive)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        string
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user retrieval",
			userID: "test-user-id",
			expectedUser: &model.User{
				ID:       "test-user-id",
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: true,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "test-user-id").
					Return(&model.User{
						ID:       "test-user-id",
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
					}, nil)
			},
		},
		{
			name:          "user not found",
			userID:        "non-existent-id",
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "non-existent-id").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			result, err := service.GetUserByID(context.Background(), tt.userID)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.ID, result.ID)
				assert.Equal(t, tt.expectedUser.Username, result.Username)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
			}
		})
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		email         string
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:  "successful user retrieval by email",
			email: "test@example.com",
			expectedUser: &model.User{
				ID:       "test-user-id",
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: true,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
					Return(&model.User{
						ID:       "test-user-id",
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
					}, nil)
			},
		},
		{
			name:          "user not found by email",
			email:         "nonexistent@example.com",
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByEmail(gomock.Any(), "nonexistent@example.com").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			result, err := service.GetUserByEmail(context.Background(), tt.email)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.Username, result.Username)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        string
		req           *dto.UpdateUserRequest
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user update",
			userID: "test-user-id",
			req: &dto.UpdateUserRequest{
				FirstName: strPtr("Updated"),
				LastName:  strPtr("Name"),
			},
			expectedUser: &model.User{
				ID:        "test-user-id",
				Username:  "testuser",
				Email:     "test@example.com",
				FirstName: strPtr("Updated"),
				LastName:  strPtr("Name"),
				IsActive:  true,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "test-user-id").
					Return(&model.User{
						ID:       "test-user-id",
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
					}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&model.User{
					ID:        "test-user-id",
					Username:  "testuser",
					Email:     "test@example.com",
					FirstName: strPtr("Updated"),
					LastName:  strPtr("Name"),
					IsActive:  true,
				}, nil)
			},
		},
		{
			name:   "user not found",
			userID: "non-existent-id",
			req: &dto.UpdateUserRequest{
				FirstName: strPtr("Updated"),
			},
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "non-existent-id").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			result, err := service.UpdateUser(context.Background(), tt.userID, tt.req)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.ID, result.ID)
				assert.Equal(t, tt.expectedUser.Username, result.Username)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        string
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:          "successful user deletion",
			userID:        "test-user-id",
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().Delete(gomock.Any(), "test-user-id").Return(nil)
			},
		},
		{
			name:          "user not found for deletion",
			userID:        "non-existent-id",
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().Delete(gomock.Any(), "non-existent-id").Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			err := service.DeleteUser(context.Background(), tt.userID)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		req           *dto.UserListRequest
		expectedUsers []*model.User
		expectedTotal int64
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name: "successful user list retrieval",
			req: &dto.UserListRequest{
				Page: 1,
				Size: 10,
			},
			expectedUsers: []*model.User{
				{
					ID:       "user-1",
					Username: "user1",
					Email:    "user1@example.com",
					IsActive: true,
				},
				{
					ID:       "user-2",
					Username: "user2",
					Email:    "user2@example.com",
					IsActive: true,
				},
			},
			expectedTotal: 2,
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]*model.User{
					{
						ID:       "user-1",
						Username: "user1",
						Email:    "user1@example.com",
						IsActive: true,
					},
					{
						ID:       "user-2",
						Username: "user2",
						Email:    "user2@example.com",
						IsActive: true,
					},
				}, int64(2), nil)
			},
		},
		{
			name: "database error",
			req: &dto.UserListRequest{
				Page: 1,
				Size: 10,
			},
			expectedUsers: nil,
			expectedTotal: 0,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, int64(0), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			users, total, err := service.ListUsers(context.Background(), tt.req)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, users)
				assert.Len(t, users, len(tt.expectedUsers))
				assert.Equal(t, tt.expectedTotal, total)
			}
		})
	}
}

func TestUserService_ActivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        string
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user activation",
			userID: "test-user-id",
			expectedUser: &model.User{
				ID:       "test-user-id",
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: true,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "test-user-id").
					Return(&model.User{
						ID:       "test-user-id",
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: false,
					}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&model.User{
					ID:       "test-user-id",
					Username: "testuser",
					Email:    "test@example.com",
					IsActive: true,
				}, nil)
			},
		},
		{
			name:          "user not found for activation",
			userID:        "non-existent-id",
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "non-existent-id").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			result, err := service.ActivateUser(context.Background(), tt.userID)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.True(t, result.IsActive)
			}
		})
	}
}

func TestUserService_DeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        string
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user deactivation",
			userID: "test-user-id",
			expectedUser: &model.User{
				ID:       "test-user-id",
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: false,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "test-user-id").
					Return(&model.User{
						ID:       "test-user-id",
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
					}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&model.User{
					ID:       "test-user-id",
					Username: "testuser",
					Email:    "test@example.com",
					IsActive: false,
				}, nil)
			},
		},
		{
			name:          "user not found for deactivation",
			userID:        "non-existent-id",
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), "non-existent-id").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.NewMockUserRepository(ctrl)
			tt.setupMock(mockRepo)
			logger := zap.NewNop()
			service := NewUserService(mockRepo, logger)
			result, err := service.DeactivateUser(context.Background(), tt.userID)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.False(t, result.IsActive)
			}
		})
	}
}

// Benchmark tests
func BenchmarkUserService_CreateUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	mockRepo := mock.NewMockUserRepository(ctrl)
	logger := zap.NewNop()
	service := NewUserService(mockRepo, logger)

	user := &model.User{
		Username:     "benchmarkuser",
		Email:        "benchmark@example.com",
		PasswordHash: "hashedpassword",
	}

	mockRepo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, assert.AnError).AnyTimes()
	mockRepo.EXPECT().GetByUsername(gomock.Any(), gomock.Any()).Return(nil, assert.AnError).AnyTimes()
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(user, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateUser(context.Background(), user)
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	mockRepo := mock.NewMockUserRepository(ctrl)
	logger := zap.NewNop()
	service := NewUserService(mockRepo, logger)

	user := &model.User{
		ID:       "benchmark-user-id",
		Username: "benchmarkuser",
		Email:    "benchmark@example.com",
		IsActive: true,
	}

	mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(user, nil).AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetUserByID(context.Background(), "benchmark-user-id")
	}
}
