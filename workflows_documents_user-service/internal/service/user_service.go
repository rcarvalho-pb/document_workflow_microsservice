package service

import (
	"errors"
	"time"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/security"
)

type ErrUserService error

var (
	ErrUserAlreadyDeactivated = errors.New("user already deactivated")
	ErrUserAlreadyActivated   = errors.New("user already activated")
	ErrUserIncorrectPassword  = errors.New("incorrect password")
)

type UserService struct {
	Repository model.UserRepository
}

func NewUserService(repo model.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Save(userDTO *dto.UserDTO) (int64, error) {
	user := userDTO.ToUserModel()
	hashedPassword, err := security.Hash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	now := time.Now().Unix()
	user.CreatedAt, user.UpdatedAt, user.Active = now, now, true
	id, err := s.Repository.Save(user)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *UserService) Update(userDTO *dto.UserDTO) error {
	user, err := s.Repository.FindByID(userDTO.ID)
	if err != nil {
		return err
	}
	updatedUser := userDTO.ToUserModel()

	if updatedUser.Name != "" {
		user.Name = updatedUser.Name
	}
	if updatedUser.LastName != "" {
		user.LastName = updatedUser.LastName
	}
	if updatedUser.Email != "" {
		user.Email = updatedUser.Email
	}
	user.UpdateUserTime()
	return s.Repository.Update(user)
}

func (s *UserService) FindByID(id int64) (*dto.UserDTO, error) {
	user, err := s.Repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	userDTO := dto.FromUserModel(user)
	return userDTO, nil
}

func (s *UserService) FindByEmail(email string) (*dto.UserDTO, error) {
	user, err := s.Repository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return dto.FromUserModel(user), nil
}

func (s *UserService) FindByName(name string) ([]*dto.UserDTO, error) {
	users, err := s.Repository.FindByName(name)
	if err != nil {
		return nil, err
	}
	usersDTO := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		userDTO := dto.FromUserModel(user)
		usersDTO[i] = userDTO
	}
	return usersDTO, nil
}

func (s *UserService) DeactivateUserByID(id int64) error {
	user, err := s.Repository.FindByID(id)
	if err != nil {
		return err
	}
	if !user.Active {
		return ErrUserAlreadyDeactivated
	}
	user.Active = false
	return s.Repository.Update(user)
}

func (s *UserService) ReactivateUserByID(id int64) error {
	user, err := s.Repository.FindByID(id)
	if err != nil {
		return err
	}
	if user.Active {
		return ErrUserAlreadyActivated
	}
	user.Active = true
	return s.Repository.Update(user)
}

func (s *UserService) UpdatePassword(id int64, changePassword *dto.ChangePassword) error {
	user, err := s.Repository.FindByID(id)
	if err != nil {
		return err
	}
	if !security.CheckPassword(changePassword.Password, user.Password) {
		return ErrUserIncorrectPassword
	}
	hashedNewPassword, err := security.Hash(changePassword.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedNewPassword
	return s.Repository.Update(user)
}
