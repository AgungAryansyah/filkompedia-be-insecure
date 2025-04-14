package repository

import (
	"database/sql"
	"errors"

	"github.com/AgungAryansyah/filkompedia-be-insecure/entity"
	"github.com/AgungAryansyah/filkompedia-be-insecure/model"
	"github.com/AgungAryansyah/filkompedia-be-insecure/pkg/response"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type IUserRepository interface {
	GetUsers(users *[]entity.User, page, pageSize int) error
	GetUser(user *entity.User, userId uuid.UUID) error
	GetUserByEmail(email string) (user *entity.User, err error)
	UpdateRole(userId uuid.UUID, roleId int) error
	EditUser(edit *model.EditProfile) error
	DeleteUser(userId uuid.UUID) error
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUsers(users *[]entity.User, page, pageSize int) error {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	query := `SELECT * FROM users ORDER BY id ASC LIMIT $1 OFFSET $2`
	err := r.db.Select(users, query, pageSize, offset)
	return err
}

func (r *UserRepository) GetUser(user *entity.User, userId uuid.UUID) error {
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(user, query, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return &response.UserNotFound
	}
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (user *entity.User, err error) {
	query := `SELECT * FROM users WHERE email = $1`

	user = &entity.User{}
	err = r.db.Get(user, query, email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, &response.UserNotFound
	}

	return user, err
}

func (r *UserRepository) UpdateRole(userId uuid.UUID, roleId int) error {
	query := `UPDATE users SET role_id = $1 WHERE id = $2`

	result, err := r.db.Exec(query, roleId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &response.UserNotFound
	}

	return nil
}

func (r *UserRepository) EditUser(edit *model.EditProfile) error {
	query := `UPDATE users SET username = '` + edit.Username + `', profile_picture = '` + edit.ProfilePicture + `' WHERE id = '` + edit.Id.String() + `'`

	_, err := r.db.Exec(query)
	return err
}

func (r *UserRepository) DeleteUser(userId uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, userId)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &response.UserNotFound
	}

	return nil
}
