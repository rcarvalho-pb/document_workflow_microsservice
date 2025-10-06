package sqlite3_test

import (
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/adapter/repository/sqlite3"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *sqlite3.DB {
	conn, err := sqlx.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "should not return error openning db in memory")
	db := &sqlite3.DB{conn}
	schema := `
	CREATE TABLE tb_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		last_name TEXT,
		email TEXT UNIQUE,
		password TEXT,
		role INTEGER,
		created_at INTEGER DEFAULT (strftime('%s', 'now')),
		updated_at INTEGER DEFAULT (strftime('%s', 'now')),
		active BOOLEAN DEFAULT 1
	);`
	_, err = db.Exec(schema)
	assert.NoError(t, err, "should not return error creating schema")
	assert.NotNil(t, db, "db should not be null")
	return db
}

func TestDB_Save(t *testing.T) {
	db := setupDB(t)
	user := &model.User{
		Name:      "Ramon",
		LastName:  "Carvalho",
		Email:     "ramon@email.com",
		Password:  "123",
		Role:      model.EMPLOYEE,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Active:    true,
	}
	id, err := db.Save(user)
	assert.NoError(t, err, "should not return error to save new user")
	assert.NotZero(t, id, "should not be zero if saved successfully")
}

func TestFindUserByID(t *testing.T) {
	t.Run("should find a saved user", func(t *testing.T) {
		db := setupDB(t)
		user := &model.User{
			Name:      "Ramon",
			LastName:  "Carvalho",
			Email:     "ramon@email.com",
			Password:  "123",
			Role:      model.EMPLOYEE,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			Active:    true,
		}
		id, err := db.Save(user)
		assert.NoError(t, err, "should not return error to save new user")
		assert.NotZero(t, id, "should not be zero if saved successfully")
		user.ID = id
		findedUser, err := db.FindByID(id)
		assert.NoError(t, err)
		assert.Equal(t, user, findedUser)
	})
	t.Run("should no find unsaved user", func(t *testing.T) {
		db := setupDB(t)
		user := &model.User{
			Name:      "Ramon",
			LastName:  "Carvalho",
			Email:     "ramon@email.com",
			Password:  "123",
			Role:      model.EMPLOYEE,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			Active:    true,
		}
		id, err := db.Save(user)
		assert.NoError(t, err, "should not return error to save new user")
		assert.NotZero(t, id, "should not be zero if saved successfully")
		findedUser, err := db.FindByID(int64(166))
		assert.Error(t, err)
		assert.Nil(t, findedUser)
	})
}
