package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"BACKEND/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (u userRepository) Create(ctx context.Context, user domain.UserCreate) (int, error) {
	var createdID int

	hashedPassword := utils.HashPassword(user.Password)

	createQuery := `INSERT INTO users (email, password, first_name, last_name, age, sex) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err := u.db.QueryRowContext(ctx, createQuery, user.Email, hashedPassword, user.FirstName, user.LastName, user.Age, user.Sex).Scan(&createdID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrAlreadyExist
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (u userRepository) GetSecure(ctx context.Context, email string) (int, string, error) {
	var hashedPassword string
	var id int

	getQuery := `SELECT id, password FROM users WHERE email = $1`

	err := u.db.QueryRowContext(ctx, getQuery, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", errs.InvalidEmail
		}
		return 0, "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return id, hashedPassword, nil
}
