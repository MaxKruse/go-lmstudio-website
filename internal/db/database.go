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

	var fErr error // Variable to store the error from f()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Transaction panicked: %v", r)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction after panic: %v", rollbackErr)
			}
		} else if fErr != nil {
			log.Printf("Transaction error: %v, rolling back...", fErr)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Failed to rollback transaction: %v", rollbackErr)
			}
		} else {
			log.Println("Committing successful transaction...")
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Failed to commit transaction: %v", commitErr)
				fErr = commitErr // Return the commit error if it happens
			}
		}
	}()

	fErr = f(tx) // Execute the provided function and store its error
	return fErr  // Return the error from f() or nil if successful
}
