package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"BACKEND/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
)

type userRepo struct {
	db *sqlx.DB
}

func InitUserRepo(
	db *sqlx.DB,
) Users {
	return &userRepo{
		db,
	}
}

func (u userRepo) Create(ctx context.Context, user domain.UserCreate) (int, error) {
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

func (u userRepo) GetByID(ctx context.Context, userID int) (domain.User, error) {
	var user domain.User

	selectQuery := `SELECT email, first_name, last_name, age, sex, photo_url
		FROM users WHERE id = $1`

	err := u.db.QueryRowxContext(ctx, selectQuery, userID).Scan(&user.Email, &user.FirstName, &user.LastName,
		&user.Age, &user.Sex, &user.PhotoUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, errs.ErrNoUser
		}
		return domain.User{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	user.ID = userID

	return user, nil
}

func (u userRepo) GetSecure(ctx context.Context, email string) (int, string, error) {
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

func (u userRepo) UpdateMain(ctx context.Context, user domain.UserUpdate) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE users SET email = $1, first_name = $2, last_name = $3, age = $4, sex = $5 WHERE id = $6`

	res, err := tx.ExecContext(ctx, updateQuery, user.Email, user.FirstName, user.LastName, user.Age, user.Sex, user.ID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, 1, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return errs.ErrNoUser
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (u userRepo) UpdatePhotoUrl(ctx context.Context, userID int, newPhotoUrl null.String) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE users SET photo_url = $1 WHERE id = $2`

	res, err := tx.ExecContext(ctx, updateQuery, newPhotoUrl, userID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, 1, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return errs.ErrNoUser
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}
