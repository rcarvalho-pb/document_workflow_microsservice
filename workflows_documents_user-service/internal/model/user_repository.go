package model

type UserRepository interface {
	Save(*User) (int64, error)
	Update(*User) error
	FindByID(int64) (*User, error)
	FindByEmail(string) (*User, error)
	FindByName(string) ([]*User, error)
}
