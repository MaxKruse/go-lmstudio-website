package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db  *sqlx.DB
	err error
)

func initDb() {
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	DB_PASS := os.Getenv("POSTGRES_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	POSTGRES_DB := os.Getenv("POSTGRES_DB")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		POSTGRES_USER, DB_PASS, DB_HOST, DB_PORT, POSTGRES_DB)

	// Open the connection
	db, err = sqlx.Open("postgres", dbURL)
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

	if db == nil {
		initDb()
	}

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
