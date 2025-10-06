package dto_test

import (
	"testing"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/dto"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestUserModelToUserDTO(t *testing.T) {
	user, _ := model.NewUser("Ramon", "Carvalho", "ramon@email.com", "123", model.ADMIN)
	userDTO := dto.FromUserModel(user)
	assert.Equal(t, user.Name, userDTO.Name)
	assert.Equal(t, user.Role.String(), userDTO.Role)
}

func TestUserDTOToUserModel(t *testing.T) {
	userDTO := &dto.UserDTO{
		ID:       1,
		Name:     "Ramon",
		LastName: "Carvalho",
		Email:    "ramon@email.com",
		Password: "123",
		Role:     "employee",
	}
	user := userDTO.ToUserModel()
	assert.Equal(t, user.Role.String(), userDTO.Role)
	assert.Equal(t, user.Email, userDTO.Email)
}
