package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "hashedpassword",
				FirstName: strPtr("Test"),
				LastName:  strPtr("User"),
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
					ID:        1,
					Username:  "testuser",
					Email:     "test@example.com",
					Password:  "hashedpassword",
					FirstName: strPtr("Test"),
					LastName:  strPtr("User"),
					IsActive:  true,
					Status:    model.UserStatusActive,
				}, nil)
			},
		},
		{
			name: "database error",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpassword",
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
				assert.Equal(t, model.UserStatusActive, result.Status)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        uint
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user retrieval",
			userID: 1,
			expectedUser: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: true,
				Status:   model.UserStatusActive,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(1)).
					Return(&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
						Status:   model.UserStatusActive,
					}, nil)
			},
		},
		{
			name:          "user not found",
			userID:        999,
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(999)).
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
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: true,
				Status:   model.UserStatusActive,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").
					Return(&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
						Status:   model.UserStatusActive,
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
		userID        uint
		req           *dto.UpdateUserRequest
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user update",
			userID: 1,
			req: &dto.UpdateUserRequest{
				FirstName: strPtr("Updated"),
				LastName:  strPtr("Name"),
			},
			expectedUser: &model.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@example.com",
				FirstName: strPtr("Updated"),
				LastName:  strPtr("Name"),
				IsActive:  true,
				Status:    model.UserStatusActive,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(1)).
					Return(&model.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						FirstName: strPtr("Original"),
						LastName:  strPtr("Name"),
						IsActive:  true,
						Status:    model.UserStatusActive,
					}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(&model.User{
						ID:        1,
						Username:  "testuser",
						Email:     "test@example.com",
						FirstName: strPtr("Updated"),
						LastName:  strPtr("Name"),
						IsActive:  true,
						Status:    model.UserStatusActive,
					}, nil)
			},
		},
		{
			name:   "user not found",
			userID: 999,
			req: &dto.UpdateUserRequest{
				FirstName: strPtr("Updated"),
			},
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(999)).
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
				assert.Equal(t, tt.expectedUser.FirstName, result.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, result.LastName)
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
			name: "successful user list",
			req: &dto.UserListRequest{
				Page:  1,
				Size:  10,
				Sort:  "created_at",
				Order: "desc",
			},
			expectedUsers: []*model.User{
				{
					ID:       1,
					Username: "user1",
					Email:    "user1@example.com",
					IsActive: true,
					Status:   model.UserStatusActive,
				},
				{
					ID:       2,
					Username: "user2",
					Email:    "user2@example.com",
					IsActive: true,
					Status:   model.UserStatusActive,
				},
			},
			expectedTotal: 2,
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return([]*model.User{
						{
							ID:       1,
							Username: "user1",
							Email:    "user1@example.com",
							IsActive: true,
							Status:   model.UserStatusActive,
						},
						{
							ID:       2,
							Username: "user2",
							Email:    "user2@example.com",
							IsActive: true,
							Status:   model.UserStatusActive,
						},
					}, int64(2), nil)
			},
		},
		{
			name: "list error",
			req: &dto.UserListRequest{
				Page: 1,
				Size: 10,
			},
			expectedUsers: nil,
			expectedTotal: 0,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().List(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), assert.AnError)
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
				assert.Equal(t, len(tt.expectedUsers), len(users))
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
		userID        uint
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user activation",
			userID: 1,
			expectedUser: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: true,
				Status:   model.UserStatusActive,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(1)).
					Return(&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: false,
						Status:   model.UserStatusInactive,
					}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
						Status:   model.UserStatusActive,
					}, nil)
			},
		},
		{
			name:          "user not found",
			userID:        999,
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(999)).
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
				assert.Equal(t, model.UserStatusActive, result.Status)
			}
		})
	}
}

func TestUserService_DeactivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tests := []struct {
		name          string
		userID        uint
		expectedUser  *model.User
		expectedError bool
		setupMock     func(*mock.MockUserRepository)
	}{
		{
			name:   "successful user deactivation",
			userID: 1,
			expectedUser: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				IsActive: false,
				Status:   model.UserStatusInactive,
			},
			expectedError: false,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(1)).
					Return(&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: true,
						Status:   model.UserStatusActive,
					}, nil)
				repo.EXPECT().Update(gomock.Any(), gomock.Any()).
					Return(&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
						IsActive: false,
						Status:   model.UserStatusInactive,
					}, nil)
			},
		},
		{
			name:          "user not found",
			userID:        999,
			expectedUser:  nil,
			expectedError: true,
			setupMock: func(repo *mock.MockUserRepository) {
				repo.EXPECT().GetByID(gomock.Any(), uint(999)).
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
				assert.Equal(t, model.UserStatusInactive, result.Status)
			}
		})
	}
}

// Benchmark tests
func BenchmarkUserService_CreateUser(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	mockRepo := mock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		Return(&model.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashedpassword",
			IsActive: true,
			Status:   model.UserStatusActive,
		}, nil).AnyTimes()

	logger := zap.NewNop()
	service := NewUserService(mockRepo, logger)

	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.CreateUser(context.Background(), user)
		require.NoError(b, err)
	}
}

func BenchmarkUserService_GetUserByID(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	mockRepo := mock.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), uint(1)).
		Return(&model.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			IsActive: true,
			Status:   model.UserStatusActive,
		}, nil).AnyTimes()

	logger := zap.NewNop()
	service := NewUserService(mockRepo, logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetUserByID(context.Background(), 1)
		require.NoError(b, err)
	}
}
