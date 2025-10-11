package main

import (
	"log"

	"github.com/rcarvalho-pb/workflows-document_user-service/internal/adapter/repository/sqlite3"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/api"
	"github.com/rcarvalho-pb/workflows-document_user-service/internal/service"
)

func main() {
	db := sqlite3.ConnectoToDB(":memory:")
	err := createTable(db)
	if err != nil {
		log.Fatal(err)
	}
	userService := service.NewUserService(db)
	grpcService := api.UserGRPCServer{}
	grpcService.Run(userService)
}

func createTable(db *sqlite3.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS tb_users (
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
	_, err := db.Exec(schema)
	return err
}
