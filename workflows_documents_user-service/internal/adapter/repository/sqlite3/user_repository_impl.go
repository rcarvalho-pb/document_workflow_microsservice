package sqlite3

import (
	"context"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
)

// Save(*User) int
// Update(*User) error
// FindByID(int) (*User, error)
// FindByEmail(string) (*User, error)
// FindByName(string) ([]*User, error)
// DeactivateUserByID(int) error
// ReactivateUserByID(int) error

func (db *DB) Save(user *model.User) int {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO tb_users (:name, )`
}
