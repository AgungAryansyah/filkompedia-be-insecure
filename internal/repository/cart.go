package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

type ICartRepository interface {
	GetUserCart(carts *[]entity.Cart, user *entity.User) error
	GetCart(cart *entity.Cart, cartId uuid.UUID) error
	AddToCart(user *entity.User, book *entity.Book, amount int) error
	RemoveFromCart(cartId uuid.UUID) error
	EditCart(cart *entity.Cart, amount int) error
	DeleteCartByBook(bookId uuid.UUID) error
	DeleteUserCart(userId uuid.UUID) error
	DeleteUser(userId uuid.UUID) error
}

type CartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) ICartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) GetUserCart(carts *[]entity.Cart, user *entity.User) error {
	query := `SELECT * FROM carts WHERE user_id = $1 AND checkout_id IS NULL`
	err := r.db.Select(carts, query, user.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return &response.CartNotFound
	}
	return err
}

func (r *CartRepository) GetCart(cart *entity.Cart, cartId uuid.UUID) error {
	query := `SELECT * FROM carts WHERE id = $1`
	err := r.db.Get(cart, query, cartId)
	if errors.Is(err, sql.ErrNoRows) {
		return &response.CartNotFound
	}
	return err
}

func (r *CartRepository) AddToCart(user *entity.User, book *entity.Book, amount int) error {
	if amount < 1 {
		return errors.New("invalid amount")
	}

	var cart entity.Cart
	if r.doesCartExist(&cart, user.Id, book.Id); cart.Amount > 0 {
		return r.addAmount(cart.Id, amount)
	}

	query := `INSERT INTO carts (id, user_id, book_id, amount, checkout_id) VALUES ($1, $2, $3, $4, NULL)`
	_, err := r.db.Exec(query, uuid.New(), user.Id, book.Id, amount)
	return err
}

func (r *CartRepository) EditCart(cart *entity.Cart, amount int) error {
	query := `UPDATE carts SET amount = $1 WHERE id = $2`
	result, err := r.db.Exec(query, amount, cart.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &response.CartNotFound
	}

	return nil
}

func (r *CartRepository) RemoveFromCart(cartId uuid.UUID) error {
	query := `DELETE FROM carts WHERE id = $1 AND checkout_id IS NULL`
	result, err := r.db.Exec(query, cartId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows deleted")
	}

	return nil
}

func (r *CartRepository) doesCartExist(cart *entity.Cart, userId uuid.UUID, bookId uuid.UUID) error {
	query := `SELECT * FROM carts WHERE user_id = $1 AND book_id = $2 AND checkout_id IS NULL LIMIT 1`
	err := r.db.Get(cart, query, userId, bookId)
	return err
}

func (r *CartRepository) addAmount(cartId uuid.UUID, amount int) error {
	query := `UPDATE carts SET amount = amount + $1 WHERE id = $2`
	_, err := r.db.Exec(query, amount, cartId)
	return err
}

func (r *CartRepository) DeleteCartByBook(bookId uuid.UUID) error {
	query := `DELETE FROM carts WHERE book_id = $1 AND checkout_id IS NULL`
	result, err := r.db.Exec(query, bookId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &response.CartNotFound
	}

	return nil
}

func (r *CartRepository) DeleteUserCart(userId uuid.UUID) error {
	query := `DELETE FROM carts WHERE user_id = $1 AND checkout_id IS NULL`
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *CartRepository) DeleteUser(userId uuid.UUID) error {
	query := `UPDATE carts SET user_id = $1 WHERE user_id = $2 AND checkout_id IS NOT NULL`
	_, err := r.db.Exec(query, uuid.Nil, userId)
	return err
}
