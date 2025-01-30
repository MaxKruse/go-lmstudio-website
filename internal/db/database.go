package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
	"github.com/maxkruse/go-lmstudio-website/internal/utils"
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

func DebugData() {

	// adds 10 books with random strings for testing

	for i := 0; i < 10; i++ {

		book := entities.Book{
			Id:            rand.Int31(),
			Title:         utils.RandomString("Title", 10),
			Author:        fmt.Sprintf("Author %d", i%3),
			Description:   utils.RandomString("Description", 20),
			ImageUrl:      fmt.Sprintf("https://example.com/image%d.jpg", rand.Intn(100)),
			PublishedDate: utils.RandomDate(),
			Isbn:          utils.RandomISBN(),
			Price:         float32(rand.Float64()*20) + 1.0,
			CreatedAt:     utils.RandomTimestamp(),
			UpdatedAt:     utils.RandomTimestamp(),
			DeletedAt:     utils.RandomTimestamp(),
		}

		err := ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
			_, err := tx.NamedExec(`INSERT INTO books (title, author, description, image_url, published_date, isbn, price) VALUES (:title, :author, :description, :image_url, :published_date, :isbn, :price)`, book)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Println(err)
		}

	}
}
