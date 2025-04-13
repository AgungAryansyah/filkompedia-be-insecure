package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

type IPaymentRepository interface {
	GetPayment(paymentId uuid.UUID) (*entity.Payment, error)
	CreatePayment(payment entity.Payment) error
	UpdatePaymentStatus(statusId int, paymentId uuid.UUID) error
	CheckUserBookPurchase(userId uuid.UUID, bookId uuid.UUID) (*bool, error)
	GetPayments(page, pageSize int) ([]entity.Payment, error)
	GetPaymentByCheckout(checkoutId uuid.UUID) (*entity.Payment, error)
	GetPaymentByUser(userId uuid.UUID) (*[]entity.Payment, error)
	DeleteUser(userId uuid.UUID) error
}

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) IPaymentRepository {
	return &PaymentRepository{db}
}

func (r *PaymentRepository) GetPayment(paymentId uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	query := `SELECT * FROM payments WHERE id = $1 LIMIT 1`
	err := r.db.Get(&payment, query, paymentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.PaymentNotFound
		}
		return nil, err
	}
	return &payment, err
}

func (r *PaymentRepository) CreatePayment(payment entity.Payment) error {
	query := `
		INSERT INTO payments (id, token, user_id, checkout_id, total_price, status_id, created_at)
		VALUES (:id, :token, :user_id, :checkout_id, :total_price, :status_id, :created_at) 
	`
	_, err := r.db.NamedExec(query, payment)
	return err
}

func (r *PaymentRepository) UpdatePaymentStatus(statusId int, paymentId uuid.UUID) error {
	query := `UPDATE payments SET status_id = $1 WHERE id = $2`
	_, err := r.db.Exec(query, statusId, paymentId)
	return err
}

func (r *PaymentRepository) CheckUserBookPurchase(userId uuid.UUID, bookId uuid.UUID) (*bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM payments
			INNER JOIN checkouts ON payments.checkout_id = checkouts.id
			INNER JOIN carts ON checkouts.id = carts.checkout_id
			WHERE payments.user_id = $1 AND carts.book_id = $2 AND payments.status_id = 1
		)
	`
	err := r.db.Get(&exists, query, userId, bookId)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func (r *PaymentRepository) GetPayments(page, pageSize int) ([]entity.Payment, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	query := `SELECT * FROM payments ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	var payments []entity.Payment

	err := r.db.Select(&payments, query, pageSize, offset)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &response.PaymentNotFound
	}

	return payments, err
}

func (r *PaymentRepository) GetPaymentByCheckout(checkoutId uuid.UUID) (*entity.Payment, error) {
	var payment entity.Payment
	query := `SELECT * FROM payments WHERE checkout_id = $1 LIMIT 1`
	err := r.db.Get(&payment, query, checkoutId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.PaymentNotFound
		}
		return nil, err
	}
	return &payment, err
}

func (r *PaymentRepository) GetPaymentByUser(userId uuid.UUID) (*[]entity.Payment, error) {
	var payment []entity.Payment
	query := `SELECT * FROM payments WHERE user_id = $1`
	err := r.db.Select(&payment, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.PaymentNotFound
		}
		return nil, err
	}
	return &payment, err
}

func (r *PaymentRepository) DeleteUser(userId uuid.UUID) error {
	query := `UPDATE payments SET user_id = $1 WHERE user_id = $2`
	_, err := r.db.Exec(query, uuid.Nil, userId)
	return err
}
