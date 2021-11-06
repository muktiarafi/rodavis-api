package config

import "database/sql"

type App struct {
	*sql.DB
}
