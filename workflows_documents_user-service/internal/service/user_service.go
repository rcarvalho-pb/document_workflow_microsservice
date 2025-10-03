package service

import (
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/security"
)

type UserService struct {
	Repository model.UserRepository
}

func (s *UserService) Save(userDTO *dto.UserDTO) (int64, error) {
	var err error
	user, err := userDTO.ToUserModel()
	if err != nil {
		return 0, err
	}
	hashedPassword, err := security.Hash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hashedPassword
	id, err := s.Repository.Save(user)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *UserService) Update(id int64, userDTO *dto.UserDTO) error {
	var err error
	user, err := s.Repository.FindByID(id)
	if err != nil {
		return err
	}
	updatedUser, err := userDTO.ToUserModel()
	if err != nil {
		return err
	}
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
	user.Active = false
	return s.Repository.Update(user)
}
