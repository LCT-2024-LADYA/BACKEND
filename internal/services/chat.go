package services

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"BACKEND/internal/repository"
	"BACKEND/pkg/log"
	"context"
	"github.com/rs/zerolog"
	"time"
)

type chatService struct {
	chatRepo       repository.Chat
	converter      converters.ChatConverter
	dbResponseTime time.Duration
	logger         zerolog.Logger
}

func InitChatService(
	chatRepo repository.Chat,
	dbResponseTime time.Duration,
	logger zerolog.Logger,
) Chat {
	return &chatService{
		chatRepo:       chatRepo,
		converter:      converters.InitChatConverter(),
		dbResponseTime: dbResponseTime,
		logger:         logger,
	}
}

func (c chatService) CreateMessage(ctx context.Context, message domain.MessageCreate) (int, time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, c.dbResponseTime)
	defer cancel()

	createdID, t, err := c.chatRepo.CreateMessage(ctx, message)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return 0, time.Time{}, err
	}

	c.logger.Info().Msg(log.Normalizer(log.CreateObject, log.Message, createdID))

	return createdID, t, nil
}

func (c chatService) GetUserChats(ctx context.Context, userID int, search string) ([]dto.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, c.dbResponseTime)
	defer cancel()

	chats, err := c.chatRepo.GetUserChats(ctx, userID, search)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return []dto.Chat{}, err
	}

	c.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Chat))

	return c.converter.ChatsDomainToDTO(chats), nil
}

func (c chatService) GetTrainerChats(ctx context.Context, trainerID int, search string) ([]dto.Chat, error) {
	ctx, cancel := context.WithTimeout(ctx, c.dbResponseTime)
	defer cancel()

	chats, err := c.chatRepo.GetTrainerChats(ctx, trainerID, search)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return []dto.Chat{}, err
	}

	c.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Chat))

	return c.converter.ChatsDomainToDTO(chats), nil
}

func (c chatService) GetChatMessage(ctx context.Context, userID, trainerID, cursor int) (dto.MessagePagination, error) {
	ctx, cancel := context.WithTimeout(ctx, c.dbResponseTime)
	defer cancel()

	messages, err := c.chatRepo.GetChatMessage(ctx, userID, trainerID, cursor)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return dto.MessagePagination{}, err
	}

	c.logger.Info().Msg(log.Normalizer(log.GetObjects, log.Message))

	return c.converter.MessagePaginationDomainToDTO(messages), nil
}
