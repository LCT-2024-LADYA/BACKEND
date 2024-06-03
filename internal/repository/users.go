package repository

import (
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func InitUserRepo(
	db *sqlx.DB,
) Users {
	return &userRepository{
		db,
	}
}

func (u userRepository) CreateVK(ctx context.Context, user domain.UserCreateVK) (int, error) {
	var createdID int

	createQuery := `INSERT INTO users (vk_id, email, first_name, last_name, age, sex) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := u.db.QueryRowContext(ctx, createQuery, user.VKID, user.Email, user.FirstName, user.LastName, user.Age, user.Sex).Scan(&createdID)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (u userRepository) CheckIfExistByVKID(ctx context.Context, VKID int) (int, error) {
	var id int

	checkQuery := `SELECT id FROM users WHERE vk_id = $1`

	err := u.db.GetContext(ctx, &id, checkQuery, VKID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return id, nil
}
