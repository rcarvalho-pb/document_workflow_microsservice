package sqlite3

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbTimeout = 10 * time.Second

type DB struct {
	*sqlx.DB
}

func ConnectoToDB(connectionString string) *DB {
	conn, err := sqlx.Open("sqlite3", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{conn}
}
