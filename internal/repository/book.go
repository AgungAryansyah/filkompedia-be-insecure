package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

type IBookRepository interface {
	GetBooks(books *[]entity.Book, page, pageSize int) error
	SearchBooks(books *[]entity.Book, page, pageSize int, searchQuery string) error
	GetBook(book *entity.Book, bookId uuid.UUID) error
	CreateBook(book *entity.Book) error
	DeleteBook(bookId uuid.UUID) error
	EditBook(edit *model.EditBook) error
}

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) IBookRepository {
	return &BookRepository{db}
}

func (r *BookRepository) GetBooks(books *[]entity.Book, page, pageSize int) error {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	query := `SELECT * FROM books WHERE author != 'This book is deleted' ORDER BY release_date DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(books, query, pageSize, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return &response.BookNotFound
	}
	return err
}

func (r *BookRepository) SearchBooks(books *[]entity.Book, page, pageSize int, searchQuery string) error {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	query := `
		SELECT * FROM books 
		WHERE 
			(title ILIKE $1 OR 
			author ILIKE $1 OR 
			description ILIKE $1) AND 
			author != 'This book is deleted' 
		ORDER BY release_date DESC 
		LIMIT $2 OFFSET $3`
	searchPattern := "%" + searchQuery + "%"
	err := r.db.Select(books, query, searchPattern, pageSize, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return &response.BookNotFound
	}
	return err
}

func (r *BookRepository) GetBook(book *entity.Book, bookId uuid.UUID) error {
	query := `SELECT * FROM books WHERE id = $1 AND author != 'This book is deleted' LIMIT 1`
	err := r.db.Get(book, query, bookId)
	if errors.Is(err, sql.ErrNoRows) {
		return &response.BookNotFound
	}
	return err
}

func (r *BookRepository) CreateBook(book *entity.Book) error {
	query := `INSERT INTO books (id, title, description, author, release_date, price, introduction, image, file) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, book.Id, book.Title, book.Description, book.Author, book.ReleaseDate, book.Price, book.Introduction, book.Image, book.File)
	return err
}

func (r *BookRepository) DeleteBook(bookId uuid.UUID) error {
	query := `
		UPDATE books 
		SET title = 'Deleted Book',
		    description = 'This book is deleted',
		    introduction = 'This book is deleted',
		    image = 'This book is deleted',
		    author = 'This book is deleted' 
		WHERE id = $1
	`

	result, err := r.db.Exec(query, bookId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &response.BookNotFound
	}

	return nil
}

func (r *BookRepository) EditBook(edit *model.EditBook) error {
	query := `
		UPDATE books 
		SET title = :title,
		    description = :description,
		    introduction = :introduction,
		    image = :image,
		    author = :author,
		    release_date = :release_date,
		    price = :price
		WHERE id = :id AND author != 'This book is deleted'
	`

	_, err := r.db.NamedExec(query, edit)
	return err
}
