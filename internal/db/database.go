package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

var (
	db  *sqlx.DB
	err error
)

func init() {
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/postgres?sslmode=disable",
		DB_USER, DB_PASS, DB_HOST,
	)

	db, err = sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
}

func ExecuteTx(ctx context.Context, f func(*sqlx.Tx) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Transaction panicked: %v", r)
			if tx.Rollback() != nil {
				log.Printf("Failed to rollback transaction after panic")
			}
		} else if err = f(tx); err != nil {
			log.Printf("Transaction error: %v, rolling back...", err)
			if tx.Rollback() != nil {
				log.Printf("Failed to rollback transaction")
			}
		} else {
			log.Println("Committing successful transaction...")
			if err := tx.Commit(); err != nil {
				log.Printf("Failed to commit transaction: %v", err)
			}
		}
	}()

	return nil // The function is designed to not return until the defer block executes
}
