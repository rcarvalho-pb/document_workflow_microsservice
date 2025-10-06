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
	None Role = iota
	ADMIN
	EMPLOYEE
	MANAGER
	FINANCIAL
)

var roles = []string{"none", "admin", "employee", "manager", "finantial"}

func (r Role) String() string {
	return roles[r]
}

func ToRole(role string) Role {
	index := Role(slices.Index(roles, role))
	if index == -1 {
		panic("role inexistente")
	}
	return index
}

type User struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
	Password  string `db:"password"`
	Role      Role   `db:"role"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	Active    bool   `db:"active"`
}

func NewUser(name, lastName, email, password string, role ...Role) (*User, error) {
	normalizedEmail := normalizeInput(email)
	if err := validateEmail(normalizedEmail); err != nil {
		return nil, err
	}
	var userRole Role
	if len(role) > 0 {
		userRole = role[0]
	} else {
		userRole = EMPLOYEE
	}
	now := time.Now().Unix()
	return &User{
		Name:      normalizeInput(name),
		LastName:  normalizeInput(lastName),
		Email:     normalizedEmail,
		Password:  normalizeInput(password),
		Role:      userRole,
		CreatedAt: now,
		UpdatedAt: now,
		Active:    true,
	}, nil
}

func validateEmail(value string) error {
	re := regexp.MustCompile(`([\w])+@([\w]+)(.[\w]{2,})(.([\w]{2,}))?$`)
	if !re.MatchString(value) {
		return ErrUserInvalidEmail
	}
	return nil
}

func normalizeInput(input string) string {
	return strings.TrimSpace(input)
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
