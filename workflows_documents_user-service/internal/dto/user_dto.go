package dto

import "github.com/rcarvalho-pb/workflows-document_user-service/internal/model"

type (
	UserDTO struct {
		ID       int64  `json:"id,omitempty"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
		Role     string `json:"role"`
	}

	ChangePassword struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}
)

func (u UserDTO) ToUserModel() (*model.User, error) {
	ub := model.UserBuilder{}
	ub.
		WithName(u.Name).
		WithLastName(u.LastName).
		WithEmail(u.Email).
		WithRole(u.Role)
	if u.Password != "" {
		ub.WithPassword(u.Password)
	}
	user, err := ub.Build()
	if err != nil {
		return nil, err
	}
	return user, err
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
