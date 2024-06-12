package repository

import (
	"BACKEND/internal/models/domain"
	"BACKEND/pkg/customerr"
	"context"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v3"
	"time"
)

type chatRepo struct {
	db                 *sqlx.DB
	entitiesPerRequest int
}

func InitChatRepo(
	db *sqlx.DB,
	entitiesPerRequest int,
) Chat {
	return &chatRepo{
		db:                 db,
		entitiesPerRequest: entitiesPerRequest,
	}
}

func (c chatRepo) CreateMessage(ctx context.Context, message domain.MessageCreate) (int, time.Time, error) {
	var createdID int
	var t time.Time

	createQuery := `INSERT INTO messages (user_id, trainer_id, message, service_id, is_to_user) VALUES ($1, $2, $3, $4, $5) RETURNING id, time`

	err := c.db.QueryRowContext(ctx, createQuery, message.UserID, message.TrainerID, message.Message, message.ServiceID, message.IsToUser).Scan(&createdID, &t)
	if err != nil {
		return 0, time.Time{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return createdID, t.UTC(), nil
}

func (c chatRepo) GetChatMessage(ctx context.Context, userID, trainerID, cursor int) (domain.MessagePagination, error) {
	query := `
		SELECT m.id, m.user_id, m.trainer_id, m.message, m.is_to_user, m.time, m.service_id
		FROM messages m
		WHERE m.user_id = $1 AND m.trainer_id = $2 AND (m.id <= $3 OR $3 = 0)
		ORDER BY m.time DESC
		LIMIT $4
	`

	rows, err := c.db.QueryContext(ctx, query, userID, trainerID, cursor, c.entitiesPerRequest+1)
	if err != nil {
		return domain.MessagePagination{}, err
	}
	defer rows.Close()

	type BasePriceNull struct {
		ID    null.Int
		Name  null.String
		Price null.Int
	}

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		var serviceID null.Int

		err := rows.Scan(&msg.ID, &msg.UserID, &msg.TrainerID, &msg.Message, &msg.IsToUser, &msg.Time, &serviceID)
		if err != nil {
			return domain.MessagePagination{}, err
		}

		if serviceID.Valid {
			a := int(serviceID.Int64)
			msg.ServiceID = &a
		} else {
			msg.ServiceID = nil
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return domain.MessagePagination{}, err
	}

	var nextCursor int
	if len(messages) == c.entitiesPerRequest+1 {
		nextCursor = messages[c.entitiesPerRequest].ID
		messages = messages[:c.entitiesPerRequest]
	}

	return domain.MessagePagination{
		Messages: messages,
		Cursor:   nextCursor,
	}, nil
}

func (c chatRepo) GetUserChats(ctx context.Context, userID int) ([]domain.Chat, error) {
	query := `
		SELECT t.id, t.photo_url, t.first_name, t.last_name, m.message, m.time
		FROM messages m
		JOIN trainers t ON m.trainer_id = t.id
		WHERE m.user_id = $1
		ORDER BY m.time DESC
		LIMIT 1;
	`

	rows, err := c.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []domain.Chat
	for rows.Next() {
		var chat domain.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.PhotoUrl,
			&chat.FirstName,
			&chat.LastName,
			&chat.LastMessage,
			&chat.TimeLastMessage,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (c chatRepo) GetTrainerChats(ctx context.Context, trainerID int) ([]domain.Chat, error) {
	query := `
		SELECT u.id, u.photo_url, u.first_name, u.last_name, m.message, m.time
		FROM messages m
		JOIN users u ON m.user_id = u.id
		WHERE m.trainer_id = $1
		ORDER BY m.time DESC
		LIMIT 1;
	`

	rows, err := c.db.QueryContext(ctx, query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []domain.Chat
	for rows.Next() {
		var chat domain.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.PhotoUrl,
			&chat.FirstName,
			&chat.LastName,
			&chat.LastMessage,
			&chat.TimeLastMessage,
		)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
