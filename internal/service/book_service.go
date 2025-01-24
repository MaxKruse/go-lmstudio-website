package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/maxkruse/go-lmstudio-website/internal/db"
	requestdtos "github.com/maxkruse/go-lmstudio-website/internal/models/dtos/request_dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
)

func GetBooks() ([]entities.Book, error) {
	var books []entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		err := tx.Select(&books, `SELECT id,
    title,
    author,
    description,
    image_url,
    published_date,
    isbn,
    price,
    created_at,
    updated_at,
    COALESCE(deleted_at::TEXT, 'Not Deleted') AS deleted_at FROM books WHERE deleted_at IS NULL ORDER BY id DESC`)

		if err != nil {
			return err
		}

		return nil
	})

	return books, err
}

func GetBookById(id int32) (entities.Book, error) {
	var book entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		err := tx.QueryRow("SELECT * FROM books WHERE id = $1 AND deleted_at IS NULL", id).Scan(&book)

		if err != nil {
			return err
		}

		return nil
	})

	return book, err
}

func CreateBook(book requestdtos.CreateBookRequest) (entities.Book, error) {
	var newBook entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		row := tx.QueryRow(`INSERT INTO books (title, author, description, image_url, published_date, isbn, price) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID`, book.Title, book.Author, book.Description, book.ImageUrl, book.PublishedDate, book.Isbn, book.Price)

		return row.Scan(&newBook.Id)
	})

	return newBook, err
}

func UpdateBook(book entities.Book) error {
	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(`UPDATE books SET title = :title, author = :author, description = :description, image_url = :image_url, published_date = :published_date, isbn = :isbn, price = :price WHERE id = :id`, book)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return err
}

func DeleteBook(id int32) error {
	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		_, err := tx.Exec(`DELETE FROM books WHERE id = $1`, id)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return err
}
