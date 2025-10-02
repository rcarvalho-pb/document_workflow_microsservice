package sqlite3

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbTimeout = 10 * time.Second

type DB struct {
	db *sqlx.DB
}

func ConnectoToDB() *DB {
	conn, err := sqlx.Open("sqlite3", "db/user_db.db")
	if err != nil {
		log.Fatal(err)
	}

	return &DB{conn}
}
