package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	RoleTable           = "roles"
	SpecializationTable = "specializations"
)

type baseRepo struct {
	db    *sqlx.DB
	table string
}

func InitBaseRepo(
	db *sqlx.DB,
	table string,
) Base {
	return &baseRepo{
		db:    db,
		table: table,
	}
}

func (c baseRepo) GetTable() string {
	return c.table
}

func (c baseRepo) Create(ctx context.Context, base domain.BaseBase) (int, error) {
	var createdID int

	createQuery := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", c.table)

	err := c.db.QueryRowContext(ctx, createQuery, base.Name).Scan(&createdID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrAlreadyExist
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (c baseRepo) Get(ctx context.Context) ([]domain.Base, error) {
	var bases []domain.Base

	getQuery := fmt.Sprintf(`SELECT id, name FROM %s ORDER BY name ASC`, c.table)

	err := c.db.SelectContext(ctx, &bases, getQuery)
	if err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	return bases, nil
}

func (c baseRepo) Delete(ctx context.Context, baseIDs []int) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = ANY($1)`, c.table)

	res, err := tx.ExecContext(ctx, deleteQuery, pq.Array(baseIDs))
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

	if int(count) != len(baseIDs) {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, 1, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return fmt.Errorf("Один из переданных объектов не существует")
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}
