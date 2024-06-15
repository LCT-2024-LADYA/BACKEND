package repository

import (
	"BACKEND/internal/errs"
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		for _, exercise := range training.Exercises {
			exerciseValueStrings = append(exerciseValueStrings, fmt.Sprintf("($%d, $%d, $%d)", argIndex, argIndex+1, argIndex+2))
			exerciseValueArgs = append(exerciseValueArgs, createdIDs[i], exercise.ID, exercise.Step)
			argIndex += 3
		}
	}

	exerciseQuery := fmt.Sprintf("INSERT INTO trainings_exercises (training_id, exercise_id, step) VALUES %s", strings.Join(exerciseValueStrings, ","))
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

func (t trainingRepo) CreateTraining(ctx context.Context, training domain.TrainingCreate) (int, error) {
	var createdID int

	tx, err := t.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Insert training
	query := `INSERT INTO trainings (name, description, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, query, training.Name, training.Description, training.UserID).Scan(&createdID)
	if err != nil {
		tx.Rollback()
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	// Insert exercises for the training and get their IDs
	if len(training.Exercises) > 0 {
		valueStrings := make([]string, 0, len(training.Exercises))
		valueArgs := make([]interface{}, 0, len(training.Exercises)*3)
		for i, exercise := range training.Exercises {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
			valueArgs = append(valueArgs, createdID, exercise.ID, exercise.Step)
		}

		exerciseQuery := fmt.Sprintf("INSERT INTO trainings_exercises (training_id, exercise_id, step) VALUES %s", strings.Join(valueStrings, ","))
		res, err := tx.ExecContext(ctx, exerciseQuery, valueArgs...)
		if err != nil {
			tx.Rollback()
			return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
		}
		count, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
		}
		if int(count) != len(training.Exercises) {
			tx.Rollback()
			return 0, errs.ErrNoExercise
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdID, nil
}

func (t trainingRepo) CreateTrainingTrainer(ctx context.Context, training domain.TrainingCreateTrainer) (int, error) {
	var createdID int

	tx, err := t.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Insert training
	query := `INSERT INTO trainings (name, description, trainer_id, wants_public) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRowContext(ctx, query, training.Name, training.Description, training.TrainerID, training.WantsPublic).Scan(&createdID)
	if err != nil {
		tx.Rollback()
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	// Insert exercises for the training and get their IDs
	if len(training.Exercises) > 0 {
		valueStrings := make([]string, 0, len(training.Exercises))
		valueArgs := make([]interface{}, 0, len(training.Exercises)*3)
		for i, exercise := range training.Exercises {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
			valueArgs = append(valueArgs, createdID, exercise.ID, exercise.Step)
		}

		exerciseQuery := fmt.Sprintf("INSERT INTO trainings_exercises (training_id, exercise_id, step) VALUES %s", strings.Join(valueStrings, ","))
		res, err := tx.ExecContext(ctx, exerciseQuery, valueArgs...)
		if err != nil {
			tx.Rollback()
			return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
		}
		count, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
		}
		if int(count) != len(training.Exercises) {
			tx.Rollback()
			return 0, errs.ErrNoExercise
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdID, nil
}

func (t trainingRepo) SetExerciseStatus(ctx context.Context, usersTrainingsID, exerciseID int, status bool) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	// Обновить статус
	updateQuery := `UPDATE user_trainings_exercises SET status = $1 WHERE users_trainings_id = $2 AND exercise_id = $3`
	res, err := tx.ExecContext(ctx, updateQuery, status, usersTrainingsID, exerciseID)
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

func (t trainingRepo) GetTrainingCovers(ctx context.Context, search string, cursor int) (domain.TrainingCoverPagination, error) {
	query := `
		SELECT t.id, t.name, t.description, COUNT(te.exercise_id) AS exercises
		FROM trainings t
		LEFT JOIN trainings_exercises te ON t.id = te.training_id
		WHERE t.user_id IS NULL
			AND (t.name LIKE '%' || $1 || '%' OR t.description LIKE '%' || $1 || '%') AND t.id >= $2
			AND (t.trainer_id IS NULL OR (t.trainer_id IS NOT NULL AND t.is_confirm = TRUE))
		GROUP BY t.id
		ORDER BY t.id
		LIMIT $3
	`
	rows, err := t.db.QueryContext(ctx, query, search, cursor, t.entitiesPerRequest+1)
	if err != nil {
		return domain.TrainingCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	var trainings []domain.TrainingCover
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

func (t trainingRepo) GetTrainingCoversByUserID(ctx context.Context, search string, userID, cursor int) (domain.TrainingCoverPagination, error) {
	query := `
		SELECT t.id, t.name, t.description, COUNT(te.exercise_id) AS exercises
		FROM trainings t
		LEFT JOIN trainings_exercises te ON t.id = te.training_id
		WHERE t.user_id = $1 AND (t.name LIKE '%' || $2 || '%' OR t.description LIKE '%' || $2 || '%') AND t.id >= $3
		GROUP BY t.id
		ORDER BY t.id
		LIMIT $4
	`
	rows, err := t.db.QueryContext(ctx, query, userID, search, cursor, t.entitiesPerRequest+1)
	if err != nil {
		return domain.TrainingCoverPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	var trainings []domain.TrainingCover
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

func (t trainingRepo) GetTrainingCoversByTrainerID(ctx context.Context, search string, trainerID, cursor int) (domain.TrainingCoverTrainerPagination, error) {
	query := `
		SELECT t.id, t.name, t.description, COUNT(te.exercise_id), t.wants_public, t.is_confirm
		FROM trainings t
		LEFT JOIN trainings_exercises te ON t.id = te.training_id
		WHERE t.trainer_id = $1 AND (t.name LIKE '%' || $2 || '%' OR t.description LIKE '%' || $2 || '%') AND t.id >= $3
		GROUP BY t.id
		ORDER BY t.id
		LIMIT $4
	`
	rows, err := t.db.QueryContext(ctx, query, trainerID, search, cursor, t.entitiesPerRequest+1)
	if err != nil {
		return domain.TrainingCoverTrainerPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	var trainings []domain.TrainingCoverTrainer
	for rows.Next() {
		var cover domain.TrainingCoverTrainer
		err := rows.Scan(&cover.ID, &cover.Name, &cover.Description, &cover.Exercises, &cover.WantsPublic, &cover.IsConfirm)
		if err != nil {
			return domain.TrainingCoverTrainerPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
		trainings = append(trainings, cover)
	}

	if err = rows.Err(); err != nil {
		return domain.TrainingCoverTrainerPagination{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	var nextCursor int
	if len(trainings) == t.entitiesPerRequest+1 {
		nextCursor = trainings[t.entitiesPerRequest].ID
		trainings = trainings[:t.entitiesPerRequest]
	}

	return domain.TrainingCoverTrainerPagination{
		Trainings: trainings,
		Cursor:    nextCursor,
	}, nil
}

func (t trainingRepo) GetScheduleTrainings(ctx context.Context, userTrainingIDs []int) ([]domain.UserTraining, error) {
	query := `
	SELECT ut.id, t.id, t.name, t.description, ut.date, ut.time_start, ut.time_end,
	       e.id, e.name, e.muscle, e.additional_muscle, e.type, e.equipment, e.difficulty,
	       e.photos, ute.sets, ute.reps, ute.weight, ute.status, te.step
	FROM users_trainings ut
		JOIN trainings t ON t.id = ut.training_id
		JOIN user_trainings_exercises ute ON ute.users_trainings_id = ut.id
		JOIN exercises e ON e.id = ute.exercise_id
		JOIN trainings_exercises te ON te.training_id = t.id AND te.exercise_id = e.id
	WHERE ut.id = ANY($1)
	ORDER BY ut.date, ut.time_start, te.step
	`
	rows, err := t.db.QueryContext(ctx, query, pq.Array(userTrainingIDs))
	if err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	trainingMap := make(map[int]domain.UserTraining)

	var training domain.UserTraining
	for rows.Next() {
		var exercise domain.Exercise

		err := rows.Scan(&training.ID, &training.TrainingID, &training.Name, &training.Description, &training.Date, &training.TimeStart, &training.TimeEnd,
			&exercise.ID, &exercise.Name, &exercise.Muscle, &exercise.AdditionalMuscle, &exercise.Type, &exercise.Equipment, &exercise.Difficulty,
			pq.Array(&exercise.Photos), &exercise.Sets, &exercise.Reps, &exercise.Weight, &exercise.Status, &exercise.Step)
		if err != nil {
			return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		if _, exists := trainingMap[training.ID]; !exists {
			training.Exercises = []domain.Exercise{}
			trainingMap[training.ID] = training
		}

		tr := trainingMap[training.ID]
		tr.Exercises = append(trainingMap[training.ID].Exercises, exercise)
		trainingMap[training.ID] = tr
	}

	if err = rows.Err(); err != nil {
		return nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	var userTrainings []domain.UserTraining
	for _, userTraining := range trainingMap {
		userTrainings = append(userTrainings, userTraining)
	}

	return userTrainings, nil
}

func (t trainingRepo) GetTraining(ctx context.Context, trainingID int) (domain.Training, error) {
	var training domain.Training

	query := `
	SELECT t.id, t.name, t.description,
	       e.id, e.name, e.muscle, e.additional_muscle, e.type, e.equipment, e.difficulty, e.photos, te.step
	FROM trainings t
	JOIN trainings_exercises te ON t.id = te.training_id
	JOIN exercises e ON te.exercise_id = e.id
	WHERE t.id = $1
	ORDER BY te.step
	`
	rows, err := t.db.QueryContext(ctx, query, trainingID)
	if err != nil {
		return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var exercise domain.ExerciseBaseStep
		err := rows.Scan(&training.ID, &training.Name, &training.Description,
			&exercise.ID, &exercise.Name, &exercise.Muscle, &exercise.AdditionalMuscle, &exercise.Type,
			&exercise.Equipment, &exercise.Difficulty, pq.Array(&exercise.Photos), &exercise.Step,
		)
		if err != nil {
			return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		training.Exercises = append(training.Exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return domain.Training{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	if len(training.Exercises) == 0 {
		return domain.Training{}, errs.ErrNoTraining
	}

	return training, nil
}

func (t trainingRepo) GetTrainingTrainer(ctx context.Context, trainingID int) (domain.TrainingTrainer, error) {
	var training domain.TrainingTrainer

	query := `
	SELECT t.id, t.name, t.description, t.wants_public, t.is_confirm,
	       e.id, e.name, e.muscle, e.additional_muscle, e.type, e.equipment, e.difficulty, e.photos, te.step
	FROM trainings t
	JOIN trainings_exercises te ON t.id = te.training_id
	JOIN exercises e ON te.exercise_id = e.id
	WHERE t.id = $1
	ORDER BY te.step
	`
	rows, err := t.db.QueryContext(ctx, query, trainingID)
	if err != nil {
		return domain.TrainingTrainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}
	defer rows.Close()

	for rows.Next() {
		var exercise domain.ExerciseBaseStep
		err := rows.Scan(&training.ID, &training.Name, &training.Description, &training.WantsPublic, &training.IsConfirm,
			&exercise.ID, &exercise.Name, &exercise.Muscle, &exercise.AdditionalMuscle, &exercise.Type,
			&exercise.Equipment, &exercise.Difficulty, pq.Array(&exercise.Photos), &exercise.Step,
		)
		if err != nil {
			return domain.TrainingTrainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		training.Exercises = append(training.Exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return domain.TrainingTrainer{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	if len(training.Exercises) == 0 {
		return domain.TrainingTrainer{}, errs.ErrNoTraining
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

	// Вставка записей в таблицу user_trainings_exercises
	if len(training.Exercises) > 0 {
		valueStrings := make([]string, 0, len(training.Exercises))
		valueArgs := make([]interface{}, 0, len(training.Exercises)*5)
		for i, exercise := range training.Exercises {
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
			valueArgs = append(valueArgs, userTrainingID, exercise.ExerciseID, exercise.Sets, exercise.Reps, exercise.Weight)
		}

		userTrainingExerciseQuery := fmt.Sprintf("INSERT INTO user_trainings_exercises (users_trainings_id, exercise_id, sets, reps, weight) VALUES %s RETURNING id", strings.Join(valueStrings, ","))
		rows, err := tx.QueryContext(ctx, userTrainingExerciseQuery, valueArgs...)
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

func (t trainingRepo) DeleteUserTraining(ctx context.Context, trainingID int) error {
	query := `DELETE FROM trainings WHERE id = $1`
	_, err := t.db.ExecContext(ctx, query, trainingID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return nil
}

func (t trainingRepo) DeleteScheduledTraining(ctx context.Context, userTrainingID int) error {
	query := `DELETE FROM users_trainings WHERE id = $1 `
	_, err := t.db.ExecContext(ctx, query, userTrainingID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return nil
}
