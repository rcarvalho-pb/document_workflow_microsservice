package sqlite3

import (
	"context"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
)

func (db *DB) Save(user *model.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `INSERT INTO tb_users
		(name, last_name, email, password, role, created_at, updated_at, active)
		VALUES
		(:name, :last_name, :email, :password, :role, :created_at, :updated_at, :active)`

	result, err := db.db.NamedExecContext(ctx, stmt, user)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (db *DB) Update(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
	UPDATE tb_users SET
	name = :name, last_name = :last_name, email = :email, password = :password, role = :role, updated_at = :updated_at, active = :active
	WHERE id = :id`

	if _, err := db.db.NamedExecContext(ctx, stmt, user); err != nil {
		return err
	}
	return nil
}

func (db *DB) FindByID(id int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM tb_users WHERE id = ?`

	user := new(model.User)
	if err := db.db.GetContext(ctx, user, query, id); err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) FindByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM tb_users WHERE email = ?`

	user := new(model.User)
	if err := db.db.GetContext(ctx, user, query, email); err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) FindByName(name string) ([]*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT * FROM tb_users WHERE name LIKE '%' || ? || '%' OR last_name LIKE '%' || ? || '%'`

	users := make([]*model.User, 0)

	if err := db.db.SelectContext(ctx, users, query, name); err != nil {
		return nil, err
	}

	return users, nil
}
