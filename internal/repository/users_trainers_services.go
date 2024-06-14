package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var UsersTrainersServicesField = map[int]string{
	1: "is_payed",
	2: "trainer_confirm",
	3: "user_confirm",
}

type usersTrainersServicesRepo struct {
	db                 *sqlx.DB
	entitiesPerRequest int
}

func InitUserTrainerServicesRepo(
	db *sqlx.DB,
	entitiesPerRequest int,
) UsersTrainersServices {
	return &usersTrainersServicesRepo{
		db:                 db,
		entitiesPerRequest: entitiesPerRequest,
	}
}

func (s usersTrainersServicesRepo) Create(ctx context.Context, service domain.UserTrainerServiceCreate) (int, error) {
	var createdID int

	createQuery := `INSERT INTO users_trainers_services (user_id, trainer_id, service_id) VALUES ($1, $2, $3) RETURNING id`

	err := s.db.QueryRowContext(ctx, createQuery, service.UserID, service.TrainerID, service.ServiceID).Scan(&createdID)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (s usersTrainersServicesRepo) Schedule(ctx context.Context, schedule domain.ScheduleService) (int, error) {
	var createdID int

	createQuery := `INSERT INTO trainer_users_trainers_services (users_trainers_services_id, date, time_start, time_end) VALUES ($1, $2, $3, $4) RETURNING id`

	err := s.db.QueryRowContext(ctx, createQuery, schedule.ScheduleID, schedule.Date, schedule.TimeStart, schedule.TimeEnd).Scan(&createdID)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, nil
}

func (s usersTrainersServicesRepo) GetSchedule(ctx context.Context, month, trainerID int) ([]domain.Schedule, error) {
	var schedules []domain.Schedule

	query := `
		SELECT date, array_agg(tuts.id) AS training_ids
		FROM trainer_users_trainers_services tuts
		JOIN users_trainers_services uts ON tuts.users_trainers_services_id = uts.id
		WHERE EXTRACT(MONTH FROM date) = $1 AND uts.trainer_id = $2
		GROUP BY date
		ORDER BY date
	`
	rows, err := s.db.QueryContext(ctx, query, month, trainerID)
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

func (s usersTrainersServicesRepo) DeleteScheduled(ctx context.Context, scheduleID int) error {
	query := `DELETE FROM trainer_users_trainers_services WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, scheduleID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return nil
}

func (s usersTrainersServicesRepo) GetUserServices(ctx context.Context, trainerID, cursor int) (domain.ServiceUserPagination, error) {

	query := `
	SELECT uts.id, uts.user_id, uts.trainer_id, uts.service_id, uts.is_payed, uts.trainer_confirm, uts.user_confirm,
	       s.id, s.name, s.price, u.id, u.first_name, u.last_name, u.age, u.sex, u.photo_url
	FROM users_trainers_services uts
		JOIN users u ON uts.user_id = u.id
		JOIN services s ON uts.service_id = s.id
	WHERE uts.trainer_id = $1 AND uts.id >= $2 LIMIT $3`

	rows, err := s.db.QueryContext(ctx, query, trainerID, cursor, s.entitiesPerRequest+1)
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
	if len(services) == s.entitiesPerRequest+1 {
		nextCursor = services[s.entitiesPerRequest].ID
		services = services[:s.entitiesPerRequest]
	}

	return domain.ServiceUserPagination{
		Services: services,
		Cursor:   nextCursor,
	}, nil
}

func (s usersTrainersServicesRepo) GetSchedulesByIDs(ctx context.Context, scheduleIDs []int) ([]domain.ScheduleServiceUser, error) {
	query := `
	SELECT uts.id, uts.user_id, uts.trainer_id, uts.service_id, uts.is_payed, uts.trainer_confirm, uts.user_confirm,
	       s.id, s.name, s.price, u.id, u.first_name, u.last_name, u.age, u.sex, u.photo_url,
	       tuts.id, tuts.date, tuts.time_start, tuts.time_end
	FROM trainer_users_trainers_services tuts
		JOIN users_trainers_services uts ON tuts.users_trainers_services_id = uts.id
		JOIN users u ON uts.user_id = u.id
		JOIN services s ON uts.service_id = s.id
	WHERE tuts.id = ANY($1)`

	rows, err := s.db.QueryContext(ctx, query, pq.Array(scheduleIDs))
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

func (s usersTrainersServicesRepo) GetTrainerServices(ctx context.Context, userID, cursor int) (domain.ServiceTrainerPagination, error) {
	query := `
	SELECT uts.id, uts.user_id, uts.trainer_id, uts.service_id, uts.is_payed, uts.trainer_confirm, uts.user_confirm,
		s.id, s.name, s.price, t.id, t.first_name, t.last_name, t.age, t.sex, t.experience, t.quote, t.photo_url,
		jsonb_agg(DISTINCT jsonb_build_object('id', r.id, 'name', r.name)) FILTER (WHERE r.id IS NOT NULL AND r.name IS NOT NULL) AS roles,
		jsonb_agg(DISTINCT jsonb_build_object('id', sp.id, 'name', sp.name)) FILTER (WHERE sp.id IS NOT NULL AND sp.name IS NOT NULL) AS specializations
	FROM users_trainers_services uts
		JOIN trainers t ON uts.trainer_id = t.id
		JOIN services s ON uts.service_id = s.id
		LEFT JOIN trainers_roles tr ON t.id = tr.trainer_id
		LEFT JOIN roles r ON tr.role_id = r.id
		LEFT JOIN trainers_specializations ts ON t.id = ts.trainer_id
		LEFT JOIN specializations sp ON ts.specialization_id = sp.id
	WHERE uts.user_id = $1 AND uts.id >= $2
	GROUP BY uts.id, s.id, t.id
	LIMIT $3`

	rows, err := s.db.QueryContext(ctx, query, userID, cursor, s.entitiesPerRequest+1)
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
	if len(services) == s.entitiesPerRequest+1 {
		nextCursor = services[s.entitiesPerRequest].ID
		services = services[:s.entitiesPerRequest]
	}

	return domain.ServiceTrainerPagination{
		Services: services,
		Cursor:   nextCursor,
	}, nil
}

func (s usersTrainersServicesRepo) UpdateStatus(ctx context.Context, field string, serviceID int, status bool) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	updateQuery := fmt.Sprintf("UPDATE users_trainers_services SET %s = $1 WHERE id = $2", field)

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

func (s usersTrainersServicesRepo) Delete(ctx context.Context, serviceID int) error {
	query := `DELETE FROM users_trainers_services WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, serviceID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return nil
}
