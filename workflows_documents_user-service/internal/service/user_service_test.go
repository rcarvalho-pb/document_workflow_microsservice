package service_test

import (
	"testing"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
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
		return args.Error(1)
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
	now := time.Now().Unix()
	repository := new(UserRepositoryMock)
	var expectedID int64 = 1
	user := &model.User{
		Name:      "Ramon",
		LastName:  "Carvalho",
		Email:     "ramon@email.com",
		Password:  "123",
		Role:      2,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	repository.On("Save", mock.AnythingOfType("*model.User")).Return(expectedID, nil)
	svc := service.NewUserService(repository)
	id, err := svc.Save(dto.FromUserModel(user))
	assert.Nil(t, err)
	assert.Equal(t, expectedID, id)
}

func TestUpdateUser(t *testing.T) {
	now := time.Now().Unix()
	user := &model.User{
		ID:        1,
		Name:      "Ramon",
		LastName:  "Carvalho",
		Email:     "ramon@email.com",
		Password:  "123",
		Role:      1,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	// updatedTime := time.Now().Unix()
	userDTO := &dto.UserDTO{
		Name:     "Emilly",
		LastName: "Coeli",
		Email:    "emilly@email.com",
	}
	mockRepo := new(UserRepositoryMock)
	mockRepo.On("FindByID", int64(1)).Return(user, nil)
	mockRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return u.Name == userDTO.Name && u.LastName == userDTO.LastName && u.Email == userDTO.Email
	})).Return(nil)
	svc := service.NewUserService(mockRepo)
	err := svc.Update(int64(1), userDTO)
	assert.Nil(t, err)
}
