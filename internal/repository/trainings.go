package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v3"
	"strings"
)

type trainingRepo struct {
	db                 *sqlx.DB
	entitiesPerRequest int
}

func InitTrainingRepo(
	db *sqlx.DB,
	entitiesPerRequest int,
) Trainings {
	return &trainingRepo{
		db:                 db,
		entitiesPerRequest: entitiesPerRequest,
	}
}

func (t trainingRepo) CreateExercises(ctx context.Context, exercises []domain.ExerciseCreateBase) ([]int, error) {
	var createdIDs []int

	tx, err := t.db.Beginx()
	if err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	valueStrings := make([]string, 0, len(exercises))
	valueArgs := make([]interface{}, 0, len(exercises)*7)
	for i, exercise := range exercises {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
		valueArgs = append(valueArgs, exercise.Name, exercise.Muscle, exercise.AdditionalMuscle, exercise.Type, exercise.Equipment, exercise.Difficulty, pq.Array(exercise.Photos))
	}

	query := fmt.Sprintf("INSERT INTO exercises (name, muscle, additional_muscle, type, equipment, difficulty, photos) VALUES %s RETURNING id", strings.Join(valueStrings, ","))
	err = tx.SelectContext(ctx, &createdIDs, query, valueArgs...)
	if err != nil {
		tx.Rollback()
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdIDs, nil
}

func (t trainingRepo) GetExercises(ctx context.Context, search string, cursor int) (domain.ExercisePagination, error) {
	var exercises []domain.ExerciseBase

	query := `
		SELECT id, name, muscle, additional_muscle, type, equipment, difficulty, photos
		FROM exercises
		WHERE id >= $1 AND name LIKE '%' || $2 || '%'
		ORDER BY id
		LIMIT $3
	`
	rows, err := t.db.QueryContext(ctx, query, cursor, search, t.entitiesPerRequest+1)
	if err != nil {
		return domain.ExercisePagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var exercise domain.ExerciseBase
		err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Muscle, &exercise.AdditionalMuscle,
			&exercise.Type, &exercise.Equipment, &exercise.Difficulty, pq.Array(&exercise.Photos))
		if err != nil {
			return domain.ExercisePagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return domain.ExercisePagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	var nextCursor int
	if len(exercises) == t.entitiesPerRequest+1 {
		nextCursor = exercises[t.entitiesPerRequest].ID
		exercises = exercises[:t.entitiesPerRequest]
	}

	return domain.ExercisePagination{
		Exercises: exercises,
		Cursor:    nextCursor,
	}, nil
}

func (t trainingRepo) CreateTrainingBases(ctx context.Context, trainings []domain.TrainingCreateBase) ([]int, error) {
	var createdIDs []int

	tx, err := t.db.Beginx()
	if err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Создание тренировок
	trainingValueStrings := make([]string, 0, len(trainings))
	trainingValueArgs := make([]interface{}, 0, len(trainings)*2)
	for i, training := range trainings {
		trainingValueStrings = append(trainingValueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		trainingValueArgs = append(trainingValueArgs, training.Name, training.Description)
	}

	trainingQuery := fmt.Sprintf("INSERT INTO trainings (name, description) VALUES %s RETURNING id", strings.Join(trainingValueStrings, ","))
	rows, err := tx.QueryxContext(ctx, trainingQuery, trainingValueArgs...)
	if err != nil {
		tx.Rollback()
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
		createdIDs = append(createdIDs, id)
	}

	if err = rows.Err(); err != nil {
		tx.Rollback()
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	// Создание упражнений для тренировок
	exerciseValueStrings := make([]string, 0)
	exerciseValueArgs := make([]interface{}, 0)
	argIndex := 1
	for i, training := range trainings {
		for _, exerciseID := range training.Exercises {
			exerciseValueStrings = append(exerciseValueStrings, fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1))
			exerciseValueArgs = append(exerciseValueArgs, createdIDs[i], exerciseID)
			argIndex += 2
		}
	}

	exerciseQuery := fmt.Sprintf("INSERT INTO trainings_exercises (training_id, exercise_id) VALUES %s", strings.Join(exerciseValueStrings, ","))
	res, err := tx.ExecContext(ctx, exerciseQuery, exerciseValueArgs...)
	if err != nil {
		tx.Rollback()
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	expectedCount := 0
	for _, training := range trainings {
		expectedCount += len(training.Exercises)
	}

	if int(count) != expectedCount {
		tx.Rollback()
		return nil, errs.ErrNoExercise
	}

	if err = tx.Commit(); err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdIDs, nil
}

func (t trainingRepo) CreateTraining(ctx context.Context, training domain.TrainingCreate) (int, []int, error) {
	var createdID int

	tx, err := t.db.Beginx()
	if err != nil {
		return 0, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Вставка тренировки
	query := `
		INSERT INTO trainings (name, description, user_id)
		VALUES ($1, $2, $3)
		RETURNING id`
	err = tx.QueryRowContext(ctx, query, training.Name, training.Description, training.UserID).Scan(&createdID)
	if err != nil {
		tx.Rollback()
		return 0, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	// Вставка упражнений для тренировки и получение их идентификаторов
	var trainingExerciseIDs []int
	if len(training.Exercises) > 0 {
		valueStrings := make([]string, 0, len(training.Exercises))
		valueArgs := make([]interface{}, 0, len(training.Exercises)*5)
		for i, exercise := range training.Exercises {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
			valueArgs = append(valueArgs, createdID, exercise.ID, exercise.Sets, exercise.Reps, exercise.Weight)
		}

		exerciseQuery := fmt.Sprintf("INSERT INTO trainings_exercises (training_id, exercise_id, sets, reps, weight) VALUES %s RETURNING id", strings.Join(valueStrings, ","))
		rows, err := tx.QueryContext(ctx, exerciseQuery, valueArgs...)
		if err != nil {
			tx.Rollback()
			return 0, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
		}
		defer rows.Close()

		for rows.Next() {
			var trainingExerciseID int
			if err := rows.Scan(&trainingExerciseID); err != nil {
				tx.Rollback()
				return 0, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
			}
			trainingExerciseIDs = append(trainingExerciseIDs, trainingExerciseID)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdID, trainingExerciseIDs, nil
}

func (t trainingRepo) SetExerciseStatus(ctx context.Context, usersTrainingsID, exerciseID int, status bool) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Найти trainings_exercises_id по exerciseID и usersTrainingsID
	var trainingsExercisesID int
	query := `
		SELECT te.id
		FROM trainings_exercises te
		JOIN user_trainings_exercises ute ON te.id = ute.trainings_exercises_id
		WHERE ute.users_trainings_id = $1 AND te.exercise_id = $2
	`
	err = tx.GetContext(ctx, &trainingsExercisesID, query, usersTrainingsID, exerciseID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNoExercise
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	// Обновить статус
	updateQuery := `UPDATE user_trainings_exercises SET status = $1 WHERE users_trainings_id = $2 AND trainings_exercises_id = $3`
	res, err := tx.ExecContext(ctx, updateQuery, status, usersTrainingsID, trainingsExercisesID)
	if err != nil {
		tx.Rollback()
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	count, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	if count != 1 {
		tx.Rollback()
		return errs.ErrNoExercise
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (t trainingRepo) GetTrainingCovers(ctx context.Context, search string, userID null.Int, cursor int) (domain.TrainingCoverPagination, error) {
	var trainings []domain.TrainingCover
	var query string
	var rows *sql.Rows
	var err error

	// Учитываем пагинацию в запросах
	if userID.Valid {
		query = `
			SELECT t.id, t.name, t.description, COUNT(te.exercise_id) AS exercises
			FROM trainings t
			LEFT JOIN trainings_exercises te ON t.id = te.training_id
			WHERE t.user_id = $1 AND (t.name LIKE '%' || $2 || '%' OR t.description LIKE '%' || $2 || '%') AND t.id >= $3
			GROUP BY t.id
			ORDER BY t.id
			LIMIT $4
		`
		rows, err = t.db.QueryContext(ctx, query, userID.Int64, search, cursor, t.entitiesPerRequest+1)
	} else {
		query = `
			SELECT t.id, t.name, t.description, COUNT(te.exercise_id) AS exercises
			FROM trainings t
			LEFT JOIN trainings_exercises te ON t.id = te.training_id
			WHERE t.user_id IS NULL AND (t.name LIKE '%' || $1 || '%' OR t.description LIKE '%' || $1 || '%') AND t.id >= $2
			GROUP BY t.id
			ORDER BY t.id
			LIMIT $3
		`
		rows, err = t.db.QueryContext(ctx, query, search, cursor, t.entitiesPerRequest+1)
	}

	if err != nil {
		return domain.TrainingCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var cover domain.TrainingCover
		err := rows.Scan(&cover.ID, &cover.Name, &cover.Description, &cover.Exercises)
		if err != nil {
			return domain.TrainingCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
		trainings = append(trainings, cover)
	}

	if err = rows.Err(); err != nil {
		return domain.TrainingCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	var nextCursor int
	if len(trainings) == t.entitiesPerRequest+1 {
		nextCursor = trainings[t.entitiesPerRequest].ID
		trainings = trainings[:t.entitiesPerRequest]
	}

	return domain.TrainingCoverPagination{
		Trainings: trainings,
		Cursor:    nextCursor,
	}, nil
}

func (t trainingRepo) GetTrainingsDate(ctx context.Context, userTrainingIDs []int) ([]domain.TrainingDate, error) {
	var trainingDates []domain.TrainingDate

	// Получение основной информации о тренировках и времени их проведения
	query := `
		SELECT ut.id, t.name, t.description, ut.date, ut.time_start, ut.time_end
		FROM trainings t
		JOIN users_trainings ut ON t.id = ut.training_id
		WHERE ut.id = ANY($1)
		ORDER BY ut.date, ut.time_start
	`
	rows, err := t.db.QueryContext(ctx, query, pq.Array(userTrainingIDs))
	if err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var trainingDate domain.TrainingDate
		err := rows.Scan(&trainingDate.ID, &trainingDate.Name, &trainingDate.Description,
			&trainingDate.Date, &trainingDate.TimeStart, &trainingDate.TimeEnd)
		if err != nil {
			return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		// Получение информации об упражнениях для каждой тренировки
		exerciseQuery := `
			SELECT e.id, e.name, e.muscle, e.additional_muscle, e.type, e.equipment, e.difficulty, e.photos,
			       te.sets, te.reps, te.weight, ute.status
			FROM user_trainings_exercises ute
			JOIN trainings_exercises te ON ute.trainings_exercises_id = te.id
			JOIN exercises e ON te.exercise_id = e.id
			WHERE ute.users_trainings_id = $1
		`
		exerciseRows, err := t.db.QueryContext(ctx, exerciseQuery, trainingDate.ID)
		if err != nil {
			return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
		}
		defer exerciseRows.Close()

		for exerciseRows.Next() {
			var exercise domain.Exercise
			err := exerciseRows.Scan(&exercise.ID, &exercise.Name, &exercise.Muscle, &exercise.AdditionalMuscle, &exercise.Type,
				&exercise.Equipment, &exercise.Difficulty, pq.Array(&exercise.Photos),
				&exercise.Sets, &exercise.Reps, &exercise.Weight, &exercise.Status)
			if err != nil {
				return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
			}

			trainingDate.Exercises = append(trainingDate.Exercises, exercise)
		}

		if err = exerciseRows.Err(); err != nil {
			return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
		}

		trainingDates = append(trainingDates, trainingDate)
	}

	if err = rows.Err(); err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return trainingDates, nil
}

func (t trainingRepo) GetTraining(ctx context.Context, trainingID int) (domain.Training, error) {
	var training domain.Training

	// Получение основной информации о тренировке
	query := `
		SELECT id, name, description
		FROM trainings
		WHERE id = $1
	`
	err := t.db.QueryRowContext(ctx, query, trainingID).Scan(&training.ID, &training.Name, &training.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Training{}, errs.ErrNoTraining
		}
		return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	// Получение информации об упражнениях
	exerciseQuery := `
		SELECT e.id, e.name, e.muscle, e.additional_muscle, e.type, e.equipment, e.difficulty, e.photos, te.sets, te.reps, te.weight
		FROM trainings_exercises te
		JOIN exercises e ON te.exercise_id = e.id
		WHERE te.training_id = $1
	`
	rows, err := t.db.QueryContext(ctx, exerciseQuery, trainingID)
	if err != nil {
		return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var exercise domain.Exercise
		err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Muscle, &exercise.AdditionalMuscle, &exercise.Type,
			&exercise.Equipment, &exercise.Difficulty, pq.Array(&exercise.Photos),
			&exercise.Sets, &exercise.Reps, &exercise.Weight)
		if err != nil {
			return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		exercise.Status = false

		training.Exercises = append(training.Exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return training, nil
}

func (t trainingRepo) ScheduleTraining(ctx context.Context, training domain.ScheduleTraining) (int, []int, error) {
	var userTrainingID int
	var userTrainingExerciseIDs []int

	tx, err := t.db.Beginx()
	if err != nil {
		return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Вставка записи в таблицу users_trainings
	query := `
		INSERT INTO users_trainings (user_id, training_id, date, time_start, time_end)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err = tx.QueryRowContext(ctx, query, training.UserID, training.TrainingID, training.Date, training.TimeStart, training.TimeEnd).Scan(&userTrainingID)
	if err != nil {
		tx.Rollback()
		return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	// Получение идентификаторов упражнений для тренировки
	exerciseQuery := `
		SELECT id
		FROM trainings_exercises
		WHERE training_id = $1`
	rows, err := tx.QueryContext(ctx, exerciseQuery, training.TrainingID)
	if err != nil {
		tx.Rollback()
		return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	var trainingExerciseIDs []int
	for rows.Next() {
		var trainingExerciseID int
		if err := rows.Scan(&trainingExerciseID); err != nil {
			tx.Rollback()
			return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
		trainingExerciseIDs = append(trainingExerciseIDs, trainingExerciseID)
	}

	// Вставка записей в таблицу user_trainings_exercises
	if len(trainingExerciseIDs) > 0 {
		valueStrings := make([]string, 0, len(trainingExerciseIDs))
		valueArgs := make([]interface{}, 0, len(trainingExerciseIDs)*2)
		for i, trainingExerciseID := range trainingExerciseIDs {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
			valueArgs = append(valueArgs, userTrainingID, trainingExerciseID)
		}

		userTrainingExerciseQuery := fmt.Sprintf("INSERT INTO user_trainings_exercises (users_trainings_id, trainings_exercises_id) VALUES %s RETURNING id", strings.Join(valueStrings, ","))
		rows, err = tx.QueryContext(ctx, userTrainingExerciseQuery, valueArgs...)
		if err != nil {
			tx.Rollback()
			return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
		}
		defer rows.Close()

		for rows.Next() {
			var userTrainingExerciseID int
			if err := rows.Scan(&userTrainingExerciseID); err != nil {
				tx.Rollback()
				return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
			}
			userTrainingExerciseIDs = append(userTrainingExerciseIDs, userTrainingExerciseID)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return userTrainingID, userTrainingExerciseIDs, nil
}

func (t trainingRepo) GetSchedule(ctx context.Context, month, userID int) ([]domain.Schedule, error) {
	var schedules []domain.Schedule

	query := `
		SELECT date, array_agg(id) AS training_ids
		FROM users_trainings
		WHERE EXTRACT(MONTH FROM date) = $1 AND user_id = $2
		GROUP BY date
		ORDER BY date
	`
	rows, err := t.db.QueryContext(ctx, query, month, userID)
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
