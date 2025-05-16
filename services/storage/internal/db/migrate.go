package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) error {
	migrationPath := "migrations"
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), dsn)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать миграцию: %w", err)
	}
	defer m.Close()
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("не удалось провести миграцию: %w", err)
	}
	return nil
}
