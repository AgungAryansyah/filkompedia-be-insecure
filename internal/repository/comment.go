package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

type ICommentRepository interface {
	GetComment(id uuid.UUID) (*entity.Comment, error)
	GetCommentByBook(bookId uuid.UUID) (*[]entity.Comment, error)
	CreateComment(comment *entity.Comment) error
	UpdateComment(comment *entity.Comment) error
	DeleteComment(id uuid.UUID) error
	DeleteCommentByBook(bookId uuid.UUID) error
	DeleteUser(userId uuid.UUID) error
}

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) ICommentRepository {
	return &CommentRepository{db}
}

func (r *CommentRepository) GetCommentByBook(bookId uuid.UUID) (*[]entity.Comment, error) {
	query := `SELECT * FROM comments WHERE book_id = $1 ORDER BY created_at DESC`
	var comments []entity.Comment
	err := r.db.Select(&comments, query, bookId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.CommentNotFound
		}
		return nil, err
	}
	return &comments, err
}

func (r *CommentRepository) GetComment(id uuid.UUID) (*entity.Comment, error) {
	query := `SELECT * FROM comments WHERE id = $1`
	var comment entity.Comment
	err := r.db.Get(&comment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.CommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) CreateComment(comment *entity.Comment) error {
	query := `
		INSERT INTO comments (id, user_id, book_id, comment, rating, created_at) 
		VALUES (:id, :user_id, :book_id, :comment, :rating, :created_at)
	`
	_, err := r.db.NamedExec(query, comment)
	return err
}

func (r *CommentRepository) UpdateComment(comment *entity.Comment) error {
	query := `
		UPDATE comments 
		SET 
			book_id = :book_id,
			comment = :comment, 
			rating = :rating, 
			created_at = :created_at
		WHERE id = :id AND user_id = :user_id
	`
	_, err := r.db.NamedExec(query, comment)
	return err
}

func (r *CommentRepository) DeleteComment(id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CommentRepository) DeleteCommentByBook(bookId uuid.UUID) error {
	query := `DELETE FROM comments WHERE book_id = $1`
	_, err := r.db.Exec(query, bookId)
	return err
}

func (r *CommentRepository) DeleteUser(userId uuid.UUID) error {
	query := `UPDATE comments SET user_id = $1 WHERE user_id = $2`
	_, err := r.db.Exec(query, uuid.Nil, userId)
	return err
}
