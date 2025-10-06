package dto

import "github.com/rcarvalho-pb/workflows-document_user-service/internal/model"

type (
	UserDTO struct {
		ID       int64  `json:"id,omitempty"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Role     string `json:"role,omitempty"`
	}

	ChangePassword struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
)

func (u UserDTO) ToUserModel() *model.User {
	user := &model.User{
		Name:     u.Name,
		LastName: u.LastName,
		Email:    u.Email,
		Password: u.Password,
	}
	if u.Password != "" {
		user.Password = u.Password
	}
	if u.Role != "" {
		user.Role = model.ToRole(u.Role)
	}
	return user
}

func FromUserModel(user *model.User) *UserDTO {
	userDTO := &UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role.String(),
	}
	return userDTO
}
