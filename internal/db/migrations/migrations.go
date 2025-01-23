package migrations

import (
	"fmt"
	"log"
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
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	log.Println("Trying to migrate up")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

	migrator, err := migrate.New("file://./migrations/", dbURL)

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
	// get db info from env
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	log.Println("Trying to migrate down")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME)

	migrator, err := migrate.New("file://./migrations/", dbURL)

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
