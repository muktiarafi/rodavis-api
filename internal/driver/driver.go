package driver

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

const (
	maxOpenDBConn = 10
	maxIdleDBConn = 5
	maxDBLifetime = 5 * time.Minute
)

func ConnectSQL(dsn string, withMigrate bool) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetMaxIdleConns(maxIdleDBConn)
	db.SetConnMaxLifetime(maxDBLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if withMigrate {
		migrationFilePath := filepath.Join(pwd, "db", "migrations")
		if err := Migration(migrationFilePath, db); err != nil {
			return nil, err
		}
	}

	return &DB{db}, nil
}
