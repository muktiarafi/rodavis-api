package config

import (
	"fmt"
	"os"
)

func CloudSQLConnection() string {
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	dbName := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	socketDir := "/cloudsql"
	return fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", user, password, dbName, socketDir, instanceConnectionName)
}
