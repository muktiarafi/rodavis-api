package driver

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.com/harta-tahta-coursera/rodavis-api/internal/logger"
)

func Migration(migrationFilePath string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationFilePath),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		logger.Notice("Migration", err.Error())
	}

	return nil
}
