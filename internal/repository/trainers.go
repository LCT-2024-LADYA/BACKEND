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
)

type trainerRepository struct {
	db *sqlx.DB
}

func InitTrainerRepo(
	db *sqlx.DB,
) Trainers {
	return &trainerRepository{
		db,
	}
}

func (t trainerRepository) Create(ctx context.Context, trainer domain.TrainerCreate) (int, error) {
	var createdID int

	hashedPassword := utils.HashPassword(trainer.Password)

	createQuery := ` INSERT INTO trainers (email, password, first_name, last_name, age, sex, experience, quote)
 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := t.db.QueryRowContext(ctx, createQuery, trainer.Email, hashedPassword, trainer.FirstName, trainer.LastName,
		trainer.Age, trainer.Sex, trainer.Experience, trainer.Quote).Scan(&createdID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrAlreadyExist
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (t trainerRepository) GetSecure(ctx context.Context, email string) (int, string, error) {
	var hashedPassword string
	var id int

	getQuery := `SELECT id, password FROM trainers WHERE email = $1`

	err := t.db.QueryRowContext(ctx, getQuery, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", errs.InvalidEmail
		}
		return 0, "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return id, hashedPassword, nil
}

func (t trainerRepository) UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := "UPDATE trainers SET first_name = $1, last_name = $2, age = $3, sex = $4, experience = $5, quote = $6 WHERE id = $7"

	res, err := tx.ExecContext(ctx, updateQuery, trainer.FirstName, trainer.LastName, trainer.Age, trainer.Sex, trainer.Experience, trainer.Quote, trainer.ID)
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
		return errs.ErrNoTrainer
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainerRepository) GetByID(ctx context.Context, trainerID int) {}
