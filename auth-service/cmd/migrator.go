package main

import (
	"errors"
	"fmt"
	"github.com/frog-in-fog/delivery-system/auth-service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *config.Config) error {
	storagePath := cfg.StoragePath
	migrationsPath := cfg.MigrationsPath
	migrationsTable := cfg.MigrationsTable

	if storagePath == "" {
		// При необходимости, можете выбрать более подходящий вариант.
		// Меня паника пока устраивает, поскольку это вспомогательная утилита.
		return errors.New("storage-path is required")
	}
	if migrationsPath == "" {
		return errors.New("migrations-path is required")
	}

	// Создаем объект мигратора, передав креды нашей БД
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		return err
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return err
		}

		return err
	}

	return nil
}
