package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать подключение к БД: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось выполнить ping: %w", err)
	}
	return db, nil
}
