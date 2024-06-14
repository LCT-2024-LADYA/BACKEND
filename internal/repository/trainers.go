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
	db                 *sqlx.DB
	entitiesPerRequest int
}

func InitTrainerRepo(
	db *sqlx.DB,
	entitiesPerRequest int,
) Trainers {
	return &trainerRepo{
		db:                 db,
		entitiesPerRequest: entitiesPerRequest,
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

	selectQuery := `
	SELECT email, first_name, last_name, age, sex, experience, quote, photo_url,
		jsonb_agg(DISTINCT jsonb_build_object('id', r.id, 'name', r.name)) FILTER (WHERE r.id IS NOT NULL AND r.name IS NOT NULL) AS roles,
		jsonb_agg(DISTINCT jsonb_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL AND t.name IS NOT NULL) AS specializations,
		jsonb_agg(DISTINCT jsonb_build_object('id', serv.id, 'name', serv.name, 'price', serv.price, 'profile_access', serv.profile_access)) FILTER (WHERE serv.id IS NOT NULL AND serv.name IS NOT NULL) AS services,
		jsonb_agg(DISTINCT jsonb_build_object('id', a.id, 'name', a.name, 'is_confirmed', a.is_confirmed)) FILTER (WHERE a.id IS NOT NULL AND a.name IS NOT NULL) AS achievements
	FROM trainers t
		LEFT JOIN trainers_roles tr ON t.id = tr.trainer_id
		LEFT JOIN roles r ON tr.role_id = r.id
		LEFT JOIN trainers_specializations ts ON t.id = ts.trainer_id
		LEFT JOIN specializations t ON ts.specialization_id = t.id
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

func (t trainerRepo) GetCovers(ctx context.Context, filters domain.FiltersTrainerCovers) (domain.TrainerCoverPagination, error) {
	selectQuery := `
	SELECT t.id, t.first_name, t.last_name, t.age, t.sex, t.experience, t.quote, t.photo_url,
		jsonb_agg(DISTINCT jsonb_build_object('id', r.id, 'name', r.name)) FILTER (WHERE r.id IS NOT NULL AND r.name IS NOT NULL) AS roles,
		jsonb_agg(DISTINCT jsonb_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL AND t.name IS NOT NULL) AS specializations
	FROM trainers t
	LEFT JOIN trainers_roles tr ON t.id = tr.trainer_id
	LEFT JOIN roles r ON tr.role_id = r.id
	LEFT JOIN trainers_specializations ts ON t.id = ts.trainer_id
	LEFT JOIN specializations t ON ts.specialization_id = t.id
	WHERE t.id IN (
		SELECT t.id
		FROM trainers t
			LEFT JOIN trainers_roles tr ON t.id = tr.trainer_id
			LEFT JOIN trainers_specializations ts ON t.id = ts.trainer_id
		WHERE t.id >= $1
			AND ($2::text IS NULL OR t.first_name LIKE '%' || $2 || '%' OR t.last_name LIKE '%' || $2 || '%')
			AND ($3::int[] IS NULL OR tr.role_id = ANY($3))
			AND ($4::int[] IS NULL OR ts.specialization_id = ANY($4))
		ORDER BY t.id)
	GROUP BY t.id
	ORDER BY t.id
	LIMIT $5`

	var trainers []domain.TrainerCover
	rows, err := t.db.QueryContext(ctx, selectQuery, filters.Cursor, filters.Search, pq.Array(filters.RoleIDs), pq.Array(filters.SpecializationIDs), t.entitiesPerRequest+1)
	if err != nil {
		return domain.TrainerCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var trainer domain.TrainerCover
		var roles, specializations []byte

		err := rows.Scan(&trainer.ID, &trainer.FirstName, &trainer.LastName, &trainer.Age, &trainer.Sex, &trainer.Experience, &trainer.Quote, &trainer.PhotoUrl, &roles, &specializations)
		if err != nil {
			return domain.TrainerCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		if len(roles) > 0 {
			if err := json.Unmarshal(roles, &trainer.Roles); err != nil {
				return domain.TrainerCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
			}
		}
		if len(specializations) > 0 {
			if err := json.Unmarshal(specializations, &trainer.Specializations); err != nil {
				return domain.TrainerCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
			}
		}

		trainers = append(trainers, trainer)
	}

	var nextCursor int
	if len(trainers) == t.entitiesPerRequest+1 {
		nextCursor = trainers[t.entitiesPerRequest].ID
		trainers = trainers[:t.entitiesPerRequest]
	}

	return domain.TrainerCoverPagination{
		Trainers: trainers,
		Cursor:   nextCursor,
	}, nil
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

		addQuery := fmt.Sprintf("INSERT INTO trainers_roles (trainer_id, role_id) VALUES %t", strings.Join(valueStrings, ","))
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

		addQuery := fmt.Sprintf("INSERT INTO trainers_specializations (trainer_id, specialization_id) VALUES %t", strings.Join(valueStrings, ","))
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

func (t trainerRepo) CreateService(ctx context.Context, service domain.ServiceCreate) (int, error) {
	var createdID int

	query := `INSERT INTO services (trainer_id, name, price, profile_access) VALUES ($1, $2, $3, $4) RETURNING id`

	err := t.db.QueryRowContext(ctx, query, service.TrainerID, service.Name, service.Price, service.ProfileAccess).Scan(&createdID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrAlreadyExist
		}
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (t trainerRepo) UpdateService(ctx context.Context, service domain.ServiceUpdate) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := `UPDATE services SET name = $1, price = $2, profile_access = $3 WHERE id = $4`

	res, err := tx.ExecContext(ctx, updateQuery, service.Name, service.Price, service.ProfileAccess, service.ID)
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

// Trainers Users Services

func (t trainerRepo) CreateServiceUser(ctx context.Context, service domain.UserTrainerServiceCreate) (int, error) {
	var createdID int

	createQuery := `INSERT INTO users_trainers_services (user_id, trainer_id, service_id) VALUES ($1, $2, $3) RETURNING id`

	err := t.db.QueryRowContext(ctx, createQuery, service.UserID, service.TrainerID, service.ServiceID).Scan(&createdID)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (t trainerRepo) Schedule(ctx context.Context, schedule domain.ScheduleService) (int, error) {
	var createdID int

	createQuery := `INSERT INTO trainer_users_trainers_services (users_trainers_services_id, date, time_start, time_end) VALUES ($1, $2, $3, $4) RETURNING id`

	err := t.db.QueryRowContext(ctx, createQuery, schedule.ScheduleID, schedule.Date, schedule.TimeStart, schedule.TimeEnd).Scan(&createdID)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (t trainerRepo) GetSchedule(ctx context.Context, month int) ([]domain.Schedule, error) {
	var schedules []domain.Schedule

	query := `
		SELECT date, array_agg(id) AS training_ids
		FROM trainer_users_trainers_services
		WHERE EXTRACT(MONTH FROM date) = $1
		GROUP BY date
		ORDER BY date
	`
	rows, err := t.db.QueryContext(ctx, query, month)
	if err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var schedule domain.Schedule
		var trainingIDs pq.Int64Array
		err := rows.Scan(&schedule.Date, &trainingIDs)
		if err != nil {
			return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		// Преобразование pq.Int64Array в []int
		schedule.TrainingIDs = make([]int, len(trainingIDs))
		for i, id := range trainingIDs {
			schedule.TrainingIDs[i] = int(id)
		}

		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return schedules, nil
}

func (t trainerRepo) DeleteScheduled(ctx context.Context, scheduleID int) error {
	query := `DELETE FROM trainer_users_trainers_services WHERE id = $1`

	_, err := t.db.ExecContext(ctx, query, scheduleID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return nil
}

func (t trainerRepo) GetUserServices(ctx context.Context, trainerID, cursor int) (domain.ServiceUserPagination, error) {

	query := `
	SELECT uts.id, uts.user_id, uts.trainer_id, uts.service_id, uts.is_payed, uts.trainer_confirm, uts.user_confirm,
	       t.id, t.name, t.price, u.id, u.first_name, u.last_name, u.age, u.sex, u.photo_url
	FROM users_trainers_services uts
		JOIN users u ON uts.user_id = u.id
		JOIN services t ON uts.service_id = t.id
	WHERE uts.trainer_id = $1 AND uts.id >= $2 LIMIT $3`

	rows, err := t.db.QueryContext(ctx, query, trainerID, cursor, t.entitiesPerRequest+1)
	if err != nil {
		return domain.ServiceUserPagination{}, err
	}
	defer rows.Close()

	var services []domain.ServiceUser
	for rows.Next() {
		var service domain.ServiceUser

		err := rows.Scan(&service.ID, &service.UserID, &service.TrainerID, &service.ServiceID, &service.IsPayed, &service.TrainerConfirm, &service.UserConfirm,
			&service.Service.ID, &service.Service.Name, &service.Service.Price, &service.User.ID, &service.User.FirstName, &service.User.LastName,
			&service.User.Age, &service.User.Sex, &service.User.PhotoUrl)
		if err != nil {
			return domain.ServiceUserPagination{}, err
		}

		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return domain.ServiceUserPagination{}, err
	}

	var nextCursor int
	if len(services) == t.entitiesPerRequest+1 {
		nextCursor = services[t.entitiesPerRequest].ID
		services = services[:t.entitiesPerRequest]
	}

	return domain.ServiceUserPagination{
		Services: services,
		Cursor:   nextCursor,
	}, nil
}

func (t trainerRepo) GetSchedulesByIDs(ctx context.Context, scheduleIDs []int) ([]domain.ScheduleServiceUser, error) {
	query := `
	SELECT uts.id, uts.user_id, uts.trainer_id, uts.service_id, uts.is_payed, uts.trainer_confirm, uts.user_confirm,
	       t.id, t.name, t.price, u.id, u.first_name, u.last_name, u.age, u.sex, u.photo_url,
	       tuts.id, tuts.date, tuts.time_start, tuts.time_end
	FROM trainer_users_trainers_services tuts
		JOIN users_trainers_services uts ON tuts.users_trainers_services_id = uts.id
		JOIN users u ON uts.user_id = u.id
		JOIN services t ON uts.service_id = t.id
	WHERE tuts.id = ANY($1)`

	rows, err := t.db.QueryContext(ctx, query, pq.Array(scheduleIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []domain.ScheduleServiceUser
	for rows.Next() {
		var service domain.ScheduleServiceUser

		err := rows.Scan(&service.ID, &service.UserID, &service.TrainerID, &service.ScheduleID, &service.IsPayed, &service.TrainerConfirm, &service.UserConfirm,
			&service.Service.ID, &service.Service.Name, &service.Service.Price, &service.User.ID, &service.User.FirstName, &service.User.LastName,
			&service.User.Age, &service.User.Sex, &service.User.PhotoUrl, &service.ScheduleID, &service.Date, &service.TimeStart, &service.TimeEnd)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}

func (t trainerRepo) GetTrainerServices(ctx context.Context, userID, cursor int) (domain.ServiceTrainerPagination, error) {
	query := `
	SELECT uts.id, uts.user_id, uts.trainer_id, uts.service_id, uts.is_payed, uts.trainer_confirm, uts.user_confirm,
		t.id, t.name, t.price, t.id, t.first_name, t.last_name, t.age, t.sex, t.experience, t.quote, t.photo_url,
		jsonb_agg(DISTINCT jsonb_build_object('id', r.id, 'name', r.name)) FILTER (WHERE r.id IS NOT NULL AND r.name IS NOT NULL) AS roles,
		jsonb_agg(DISTINCT jsonb_build_object('id', sp.id, 'name', sp.name)) FILTER (WHERE sp.id IS NOT NULL AND sp.name IS NOT NULL) AS specializations
	FROM users_trainers_services uts
		JOIN trainers t ON uts.trainer_id = t.id
		JOIN services t ON uts.service_id = t.id
		LEFT JOIN trainers_roles tr ON t.id = tr.trainer_id
		LEFT JOIN roles r ON tr.role_id = r.id
		LEFT JOIN trainers_specializations ts ON t.id = ts.trainer_id
		LEFT JOIN specializations sp ON ts.specialization_id = sp.id
	WHERE uts.user_id = $1 AND uts.id >= $2
	GROUP BY uts.id, t.id, t.id
	LIMIT $3`

	rows, err := t.db.QueryContext(ctx, query, userID, cursor, t.entitiesPerRequest+1)
	if err != nil {
		return domain.ServiceTrainerPagination{}, err
	}
	defer rows.Close()

	var services []domain.ServiceTrainer
	for rows.Next() {
		var service domain.ServiceTrainer
		var roles, specializations []byte

		err := rows.Scan(&service.ID, &service.UserID, &service.TrainerID, &service.ServiceID, &service.IsPayed, &service.TrainerConfirm, &service.UserConfirm,
			&service.Service.ID, &service.Service.Name, &service.Service.Price, &service.Trainer.ID, &service.Trainer.FirstName, &service.Trainer.LastName,
			&service.Trainer.Age, &service.Trainer.Sex, &service.Trainer.Experience, &service.Trainer.Quote, &service.Trainer.PhotoUrl, &roles, &specializations)
		if err != nil {
			return domain.ServiceTrainerPagination{}, err
		}

		if len(roles) > 0 {
			if err := json.Unmarshal(roles, &service.Trainer.Roles); err != nil {
				return domain.ServiceTrainerPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
			}
		}
		if len(specializations) > 0 {
			if err := json.Unmarshal(specializations, &service.Trainer.Specializations); err != nil {
				return domain.ServiceTrainerPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.JsonErr, Err: err})
			}
		}

		services = append(services, service)
	}

	if err = rows.Err(); err != nil {
		return domain.ServiceTrainerPagination{}, err
	}

	var nextCursor int
	if len(services) == t.entitiesPerRequest+1 {
		nextCursor = services[t.entitiesPerRequest].ID
		services = services[:t.entitiesPerRequest]
	}

	return domain.ServiceTrainerPagination{
		Services: services,
		Cursor:   nextCursor,
	}, nil
}

func (t trainerRepo) UpdateStatus(ctx context.Context, field string, serviceID int, status bool) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := fmt.Sprintf("UPDATE users_trainers_services SET %t = $1 WHERE id = $2", field)

	res, err := tx.ExecContext(ctx, updateQuery, status, serviceID)
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

func (t trainerRepo) DeleteServiceUser(ctx context.Context, serviceID int) error {
	query := `DELETE FROM users_trainers_services WHERE id = $1`

	_, err := t.db.ExecContext(ctx, query, serviceID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return nil
}
