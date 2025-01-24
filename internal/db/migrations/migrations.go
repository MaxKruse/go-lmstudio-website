package migrations

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func getMigrator() (*migrate.Migrate, error) {
	// get db info from env
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	DB_PASS := os.Getenv("POSTGRES_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	POSTGRES_DB := os.Getenv("POSTGRES_DB")

	log.Println("Trying to migrate")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		POSTGRES_USER, DB_PASS, DB_HOST, DB_PORT, POSTGRES_DB)

	return migrate.New("file://./migrations/", dbURL)
}

func MigrateUp() error {
	migrator, err := getMigrator()

	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Printf("Migration Up successful")
	return nil
}

func MigrateDown() error {
	migrator, err := getMigrator()

	if err != nil {
		return err
	}

	err = migrator.Down()
	if err != nil {
		return err
	}

	log.Printf("Migration Down successful")
	return nil
}
