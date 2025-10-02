package model

import (
	"errors"
	"regexp"
	"slices"
	"strings"
	"time"
)

type UserError error

var (
	ErrUserAlreadyDeactivate = errors.New("user is already deactivated")
	ErrUserAlreadyActive     = errors.New("user is already active")
	ErrUserInvalidEmail      = errors.New("user with invalid email")
)

type Role int

const (
	ADMIN Role = iota + 1
	EMPLOYEE
	MANAGER
	FINANCIAL
)

var roles = []string{"admin", "employee", "manager", "finantial"}

func (r Role) String() string {
	return roles[r-1]
}

func ToRole(role string) Role {
	index := Role(slices.Index(roles, role) + 1)
	if index == 0 {
		panic("role inexistente")
	}
	return index
}

type (
	userMod     func(*User)
	UserBuilder struct {
		actions []userMod
	}
	User struct {
		ID        int    `db:"id"`
		Name      string `db:"name"`
		LastName  string `db:"last_name"`
		Email     string `db:"email"`
		Password  string `db:"password"`
		Role      Role   `db:"role"`
		CreatedAt int64  `db:"created_at"`
		UpdatedAt int64  `db:"updated_at"`
		Active    bool   `db:"active"`
	}
)

func validateEmail(value string) error {
	re := regexp.MustCompile(`([\w])+@([\a]+).([\a]{2,5})([.][\a]{2,5})?`)
	if !re.MatchString(value) {
		return ErrUserInvalidEmail
	}
	return nil
}

func (u *User) NormalizeUser() {
	u.Name = strings.TrimSpace(u.Name)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(u.Email)
}

func (u *User) DeactivateUser() error {
	if u.Active {
		u.Active = false
		u.UpdateUserTime()
		return nil
	} else {
		return ErrUserAlreadyDeactivate
	}
}

func (u *User) ReactivateUser() error {
	if !u.Active {
		u.Active = true
		u.UpdateUserTime()
		return nil
	} else {
		return ErrUserAlreadyActive
	}
}

func (u *User) UpdateUserTime() {
	u.UpdatedAt = time.Now().Unix()
}

func (ub *UserBuilder) Build() (*User, error) {
	now := time.Now().Unix()
	u := &User{
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}
	for _, a := range ub.actions {
		a(u)
	}
	u.NormalizeUser()
	if err := validateEmail(u.Email); err != nil {
		return nil, err
	}
	return u, nil
}

func (b *UserBuilder) WithName(value string) *UserBuilder {
	b.actions = append(b.actions, func(u *User) {
		u.Name = value
	})
	return b
}

func (b *UserBuilder) WithLastName(value string) *UserBuilder {
	b.actions = append(b.actions, func(u *User) {
		u.LastName = value
	})
	return b
}

func (b *UserBuilder) WithEmail(value string) *UserBuilder {
	b.actions = append(b.actions, func(u *User) {
		u.Email = value
	})
	return b
}

func (b *UserBuilder) WithPassword(value string) *UserBuilder {
	b.actions = append(b.actions, func(u *User) {
		u.Password = value
	})
	return b
}

func (b *UserBuilder) WithRole(value ...any) *UserBuilder {
	b.actions = append(b.actions, func(u *User) {
		if len(value) == 0 {
			u.Role = Role(2)
		} else {
			switch v := value[0].(type) {
			case int:
				u.Role = Role(v)
			case string:
				u.Role = ToRole(v)
			default:
				panic("invalid role type")
			}
		}
	})
	return b
}
