package service_test

import (
	"testing"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/security"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) Save(user *model.User) (int64, error) {
	args := r.Called(user)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}

	return args.Get(0).(int64), args.Error(1)
}

func (r *UserRepositoryMock) Update(user *model.User) error {
	args := r.Called(user)
	if args.Get(0) == nil {
		return args.Error(0)
	}
	return nil
}

func (r *UserRepositoryMock) FindByID(id int64) (*model.User, error) {
	args := r.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *UserRepositoryMock) FindByEmail(email string) (*model.User, error) {
	args := r.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (r *UserRepositoryMock) FindByName(name string) ([]*model.User, error) {
	args := r.Called(name)
	return args.Get(0).([]*model.User), args.Error(1)
}

func TestSaveUser(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	mockRepo.AssertExpectations(t)
	userDTO := &dto.UserDTO{
		Name:     "Ramon",
		LastName: "Carvalho",
		Email:    "ramon@email.com",
		Password: "123",
		Role:     "employee",
	}
	mockRepo.On("Save", mock.AnythingOfType("*model.User")).Return(int64(1), nil)
	srv := service.NewUserService(mockRepo)
	id, err := srv.Save(userDTO)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), id)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	srv := service.NewUserService(mockRepo)
	mockRepo.AssertExpectations(t)
	now := time.Now().Unix()
	expectedUser := &model.User{
		ID:        1,
		Name:      "Ramon",
		LastName:  "Carvalho",
		Email:     "ramon@email.com",
		Password:  "123",
		Role:      model.EMPLOYEE,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	userDTO := &dto.UserDTO{
		ID:       1,
		Name:     "Emilly",
		LastName: "Coeli",
		Email:    "emilly@email",
		Password: "123",
		Role:     "employee",
	}
	mockRepo.On("FindByID", int64(1)).Return(expectedUser, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return u.Name == userDTO.Name &&
			u.LastName == userDTO.LastName &&
			u.Email == userDTO.Email
	})).Return(nil)
	err := srv.Update(userDTO)
	assert.Nil(t, err)
	mockRepo.AssertCalled(t, "FindByID", int64(1))
	mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("*model.User"))
}

func TestDeactivateUser(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	srv := service.NewUserService(mockRepo)
	mockRepo.AssertExpectations(t)
	now := time.Now().Unix()
	expectedUser := &model.User{
		ID:        1,
		Name:      "Ramon",
		LastName:  "Carvalho",
		Email:     "ramon@email.com",
		Password:  "123",
		Role:      model.EMPLOYEE,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	mockRepo.On("FindByID", int64(1)).Return(expectedUser, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return u.Active == false
	})).Return(nil)
	err := srv.DeactivateUserByID(1)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FindByID", int64(1))
	mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("*model.User"))
	mockRepo.AssertExpectations(t)
}

func TestReactivateUser(t *testing.T) {
	mockRepo := new(UserRepositoryMock)
	srv := service.NewUserService(mockRepo)
	mockRepo.AssertExpectations(t)
	now := time.Now().Unix()
	expectedUser := &model.User{
		ID:        1,
		Name:      "Ramon",
		LastName:  "Carvalho",
		Email:     "ramon@email.com",
		Password:  "123",
		Role:      model.EMPLOYEE,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    false,
	}
	mockRepo.On("FindByID", int64(1)).Return(expectedUser, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return u.Active == true
	})).Return(nil)
	err := srv.ReactivateUserByID(1)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "FindByID", int64(1))
	mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("*model.User"))
	mockRepo.AssertExpectations(t)
}

func TestUpdatePassword(t *testing.T) {
	t.Run("should update password if password is correct", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		srv := service.NewUserService(mockRepo)
		password := "123"
		hashedPassword, err := security.Hash(password)
		assert.NoError(t, err)
		changePasswordRequest := &dto.ChangePassword{
			Password:    password,
			NewPassword: "456",
		}
		user := &model.User{
			ID:        1,
			Name:      "Ramon",
			LastName:  "Carvalho",
			Email:     "ramon@email.com",
			Password:  hashedPassword,
			Role:      model.EMPLOYEE,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: (time.Now().Add(1 * time.Hour)).Unix(),
			Active:    true,
		}
		mockRepo.On("FindByID", int64(1)).Return(user, nil)
		mockRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
			return security.CheckPassword(changePasswordRequest.NewPassword, u.Password)
		})).Return(nil)
		err = srv.UpdatePassword(int64(1), changePasswordRequest)
		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "FindByID", int64(1))
		mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("*model.User"))
		mockRepo.AssertExpectations(t)
	})
	t.Run("should fail with wrong password", func(t *testing.T) {
		mockRepo := new(UserRepositoryMock)
		srv := service.NewUserService(mockRepo)
		changePasswordRequest := &dto.ChangePassword{
			Password:    "123",
			NewPassword: "456",
		}
		hashedPassword, err := security.Hash("aonde")
		assert.NoError(t, err)
		user := &model.User{
			ID:        1,
			Name:      "Ramon",
			LastName:  "",
			Email:     "",
			Password:  hashedPassword,
			Role:      0,
			CreatedAt: 0,
			UpdatedAt: 0,
			Active:    false,
		}
		mockRepo.On("FindByID", int64(1)).Return(user, nil)
		err = srv.UpdatePassword(int64(1), changePasswordRequest)
		assert.Error(t, err)
		mockRepo.AssertCalled(t, "FindByID", int64(1))
		mockRepo.AssertExpectations(t)
	})
}
