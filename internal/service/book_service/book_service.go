package book_service

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/maxkruse/go-lmstudio-website/internal/db"
	requestdtos "github.com/maxkruse/go-lmstudio-website/internal/models/dtos/request_dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
)

const query_SELECT_BOOK = `SELECT id,
    title,
    author,
    description,
    image_url,
    published_date,
    isbn,
    price,
    created_at,
    updated_at
	FROM books`

func Get() ([]entities.Book, error) {
	var books []entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		err := tx.Select(&books, fmt.Sprintf("%s WHERE deleted_at IS NULL ORDER BY id DESC", query_SELECT_BOOK))

		if err != nil {
			return err
		}

		return nil
	})

	return books, err
}

func GetById(id int32) (entities.Book, error) {
	var book entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		err := tx.Get(&book, fmt.Sprintf("%s WHERE id=$1 AND deleted_at IS NULL", query_SELECT_BOOK), id)

		if err != nil {
			return err
		}

		return nil
	})

	return book, err
}

func Create(book requestdtos.CreateBookRequest) (entities.Book, error) {
	var newBook entities.Book

	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		row := tx.QueryRow(`INSERT INTO books (title, author, description, image_url, published_date, isbn, price) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID`, book.Title, book.Author, book.Description, book.ImageUrl, book.PublishedDate, book.Isbn, book.Price)

		return row.Scan(&newBook.Id)
	})

	// get the new book data from that id
	if err != nil {
		return newBook, err
	}

	newBook, err = GetById(newBook.Id)

	return newBook, err
}

func Update(book entities.Book) error {
	err := db.ExecuteTx(context.Background(), func(tx *sqlx.Tx) error {
		_, err := tx.NamedExec(`UPDATE books SET title = :title, author = :author, description = :description, image_url = :image_url, published_date = :published_date, isbn = :isbn, price = :price WHERE id=:id`, book)

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

func Delete(id int32) error {
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

// custom funcs

func GetBooksBetweenPrice(ctx context.Context, price_min float64, price_max float64) ([]entities.Book, error) {
	var val []entities.Book

	query := fmt.Sprintf("%s WHERE price>$1 AND price<$2 AND deleted_at IS NULL ORDER BY id DESC", query_SELECT_BOOK)

	err := db.ExecuteTx(ctx, func(tx *sqlx.Tx) error {
		rows, err := tx.Queryx(query, price_min, price_max)

		if err != nil {
			return err
		}

		for rows.Next() {
			book := entities.Book{}
			err := rows.StructScan(&book)
			if err != nil {
				return err
			}

			val = append(val, book)
		}

		return nil
	})

	return val, err
}
