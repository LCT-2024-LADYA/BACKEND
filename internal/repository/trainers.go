package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"BACKEND/pkg/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
	"strings"
)

type trainerRepo struct {
	db *sqlx.DB
}

func InitTrainerRepo(
	db *sqlx.DB,
) Trainers {
	return &trainerRepo{
		db,
	}
}

func (t trainerRepo) Create(ctx context.Context, trainer domain.TrainerCreate) (int, error) {
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

func (t trainerRepo) GetSecure(ctx context.Context, email string) (int, string, error) {
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

func (t trainerRepo) GetByID(ctx context.Context, trainerID int) (domain.Trainer, error) {
	var (
		trainer                                        domain.Trainer
		roles, specializations, services, achievements []byte
	)

	selectQuery := `SELECT email, first_name, last_name, age, sex, experience, quote, photo_url,
					jsonb_agg(DISTINCT jsonb_build_object('id', r.id, 'name', r.name)) FILTER (WHERE r.id IS NOT NULL AND r.name IS NOT NULL) AS roles,
					jsonb_agg(DISTINCT jsonb_build_object('id', s.id, 'name', s.name)) FILTER (WHERE s.id IS NOT NULL AND s.name IS NOT NULL) AS specializations,
					jsonb_agg(DISTINCT jsonb_build_object('id', serv.id, 'name', serv.name, 'price', serv.price)) FILTER (WHERE serv.id IS NOT NULL AND serv.name IS NOT NULL) AS services,
					jsonb_agg(DISTINCT jsonb_build_object('id', a.id, 'name', a.name, 'is_confirmed', a.is_confirmed)) FILTER (WHERE a.id IS NOT NULL AND a.name IS NOT NULL) AS achievements
					FROM trainers t
					LEFT JOIN trainers_roles tr ON t.id = tr.trainer_id
					LEFT JOIN roles r ON tr.role_id = r.id
					LEFT JOIN trainers_specializations ts ON t.id = ts.trainer_id
					LEFT JOIN specializations s ON ts.specialization_id = s.id
					LEFT JOIN services serv ON t.id = serv.trainer_id
					LEFT JOIN achievements a ON t.id = a.trainer_id
					WHERE t.id = $1 GROUP BY t.id`

	err := t.db.QueryRowxContext(ctx, selectQuery, trainerID).Scan(&trainer.Email, &trainer.FirstName, &trainer.LastName,
		&trainer.Age, &trainer.Sex, &trainer.Experience, &trainer.Quote, &trainer.PhotoUrl, &roles, &specializations, &services, &achievements)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Trainer{}, errs.ErrNoTrainer
		}
		return domain.Trainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	trainer.ID = trainerID

	if len(roles) > 0 {
		if err := json.Unmarshal(roles, &trainer.Roles); err != nil {
			return domain.Trainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
		}
	}
	if len(specializations) > 0 {
		if err := json.Unmarshal(specializations, &trainer.Specializations); err != nil {
			return domain.Trainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
		}
	}
	if len(services) > 0 {
		if err := json.Unmarshal(services, &trainer.Services); err != nil {
			return domain.Trainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
		}
	}
	if len(achievements) > 0 {
		if err := json.Unmarshal(achievements, &trainer.Achievements); err != nil {
			return domain.Trainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
		}
	}

	return trainer, nil
}

func (t trainerRepo) UpdateMain(ctx context.Context, trainer domain.TrainerUpdate) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE trainers SET email = $1, first_name = $2, last_name = $3, age = $4, sex = $5,
		experience = $6, quote = $7 WHERE id = $8`

	res, err := tx.ExecContext(ctx, updateQuery, trainer.Email, trainer.FirstName, trainer.LastName, trainer.Age,
		trainer.Sex, trainer.Experience, trainer.Quote, trainer.ID)
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

func (t trainerRepo) UpdatePhotoUrl(ctx context.Context, trainerID int, newPhotoUrl null.String) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE trainers SET photo_url = $1 WHERE id = $2`

	res, err := tx.ExecContext(ctx, updateQuery, newPhotoUrl, trainerID)
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

func (t trainerRepo) UpdateRoles(ctx context.Context, trainerID int, roleIDs []int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	deleteQuery := `DELETE FROM trainers_roles WHERE trainer_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, trainerID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	if len(roleIDs) > 0 {
		valueStrings := make([]string, 0, len(roleIDs))
		valueArgs := make([]interface{}, 0, len(roleIDs)+1)
		valueArgs = append(valueArgs, trainerID)
		for i, id := range roleIDs {
			valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
			valueArgs = append(valueArgs, id)
		}

		addQuery := fmt.Sprintf("INSERT INTO trainers_roles (trainer_id, role_id) VALUES %s", strings.Join(valueStrings, ","))
		res, err := tx.ExecContext(ctx, addQuery, valueArgs...)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return errs.ErrAlreadyExist
			}
			if errors.As(err, &pqErr) && pqErr.Code == "23503" {
				return errs.ErrNoRole
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

		if int(count) != len(roleIDs) {
			if rbErr := tx.Rollback(); rbErr != nil {
				return customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, len(roleIDs), count)},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}
			return errs.ErrNoRole
		}
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainerRepo) UpdateSpecializations(ctx context.Context, trainerID int, specializationIDs []int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	deleteQuery := `DELETE FROM trainers_specializations WHERE trainer_id = $1`
	_, err = tx.ExecContext(ctx, deleteQuery, trainerID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	if len(specializationIDs) > 0 {
		valueStrings := make([]string, 0, len(specializationIDs))
		valueArgs := make([]interface{}, 0, len(specializationIDs)+1)
		valueArgs = append(valueArgs, trainerID)
		for i, id := range specializationIDs {
			valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
			valueArgs = append(valueArgs, id)
		}

		addQuery := fmt.Sprintf("INSERT INTO trainers_specializations (trainer_id, specialization_id) VALUES %s", strings.Join(valueStrings, ","))
		res, err := tx.ExecContext(ctx, addQuery, valueArgs...)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				return errs.ErrAlreadyExist
			}
			if errors.As(err, &pqErr) && pqErr.Code == "23503" {
				return errs.ErrNoSpecialization
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

		if int(count) != len(specializationIDs) {
			if rbErr := tx.Rollback(); rbErr != nil {
				return customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, len(specializationIDs), count)},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}
			return errs.ErrNoSpecialization
		}
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainerRepo) CreateService(ctx context.Context, trainerID int, name string, price int) (int, error) {
	var createdID int

	query := `INSERT INTO services (trainer_id, name, price) VALUES ($1, $2, $3) RETURNING id`

	err := t.db.QueryRowContext(ctx, query, trainerID, name, price).Scan(&createdID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrAlreadyExist
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (t trainerRepo) UpdateService(ctx context.Context, serviceID int, name string, price int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE services SET name = $1, price = $2 WHERE id = $3`

	res, err := tx.ExecContext(ctx, updateQuery, name, price, serviceID)
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
		return errs.ErrNoService
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainerRepo) DeleteService(ctx context.Context, trainerID, serviceID int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `DELETE FROM services WHERE id = $1 AND trainer_id = $2`

	res, err := tx.ExecContext(ctx, updateQuery, serviceID, trainerID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
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
		return errs.ErrNoService
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainerRepo) CreateAchievement(ctx context.Context, trainerID int, achievement string) (int, error) {
	var createdID int

	query := `INSERT INTO achievements (trainer_id, name) VALUES ($1, $2) RETURNING id`

	err := t.db.QueryRowContext(ctx, query, trainerID, achievement).Scan(&createdID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrAlreadyExist
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (t trainerRepo) UpdateAchievementStatus(ctx context.Context, achievementID int, status bool) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE achievements SET is_confirmed = $1 WHERE id = $2`

	res, err := tx.ExecContext(ctx, updateQuery, status, achievementID)
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
		return errs.ErrNoAchievement
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainerRepo) DeleteAchievement(ctx context.Context, trainerID, achievementID int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `DELETE FROM achievements WHERE id = $1 AND trainer_id = $2`

	res, err := tx.ExecContext(ctx, updateQuery, achievementID, trainerID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
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
		return errs.ErrNoAchievement
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}
