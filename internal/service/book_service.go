package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/maxkruse/go-lmstudio-website/internal/db"
	requestdtos "github.com/maxkruse/go-lmstudio-website/internal/models/dtos/request_dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
)

func GetBooks() []entities.Book {
	var books []entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		err := tx.Select(&books, "SELECT * FROM books WHERE deleted_at IS NULL ORDER BY id DESC")

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	// TODO: Implement the actual logic to fetch books from the database

	return books
}

func GetBookById(id int32) entities.Book {
	var book entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		err := tx.Get(&book, "SELECT * FROM books WHERE id = $1 AND deleted_at IS NULL", id)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	// TODO: Implement the actual logic to fetch a book from the database

	return book
}

func CreateBook(book requestdtos.CreateBookRequest) error {
	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(`INSERT INTO books (title, author, description, image_url, published_date, isbn, price) VALUES (:title, :author, :description, :image_url, :published_date, :isbn, :price)`, book)

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
