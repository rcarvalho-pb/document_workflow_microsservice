package model

type UserRepository interface {
	Save(*User) int
	Update(*User) error
	FindByID(int) (*User, error)
	FindByEmail(string) (*User, error)
	FindByName(string) ([]*User, error)
	DeactivateUserByID(int) error
	ReactivateUserByID(int) error
}
