package migrations

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp() error {
	// get db info from env
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/postgres?sslmode=disable", DB_USER, DB_PASS, DB_HOST)

	migrator, err := migrate.New("file://./internal/db/migrations", connectionString)

	if err != nil {
		return err
	}

	return migrator.Up()
}

func MigrateDown() error {
	// get db info from env
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_HOST := os.Getenv("DB_HOST")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/postgres?sslmode=disable", DB_USER, DB_PASS, DB_HOST)

	migrator, err := migrate.New("file://./internal/db/migrations", connectionString)

	if err != nil {
		return err
	}

	return migrator.Down()
}
