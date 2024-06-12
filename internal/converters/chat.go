package converters

import (
	"BACKEND/internal/models/domain"
	"BACKEND/internal/models/dto"
	"fmt"
	"time"
)

type ChatConverter interface {
	MessageDomainToDTO(message domain.Message) dto.Message
	MessagesDomainToDTO(messages []domain.Message) []dto.Message
	MessagePaginationDomainToDTO(pagination domain.MessagePagination) dto.MessagePagination
	ChatDomainToDTO(chat domain.Chat) dto.Chat
	ChatsDomainToDTO(chats []domain.Chat) []dto.Chat

	MessageGetToMessageCreate(message domain.MessageGet, isTrainer bool, userID int) domain.MessageCreate
	MessageCreateToMessage(message domain.MessageCreate, id int, t time.Time) domain.Message
}

type chatConverter struct {
	baseConverter BaseConverter
}

func InitChatConverter() ChatConverter {
	return &chatConverter{
		baseConverter: InitBaseConverter(),
	}
}

func (c chatConverter) MessageDomainToDTO(message domain.Message) dto.Message {
	return dto.Message{
		ID:        message.ID,
		UserID:    message.UserID,
		TrainerID: message.TrainerID,
		Message:   message.Message,
		Service:   message.ServiceID,
		IsToUser:  message.IsToUser,
		Time:      message.Time,
	}
}

func (c chatConverter) MessagesDomainToDTO(messages []domain.Message) []dto.Message {
	result := make([]dto.Message, len(messages))

	for i, message := range messages {
		result[i] = c.MessageDomainToDTO(message)
	}

	return result
}

func (c chatConverter) MessagePaginationDomainToDTO(pagination domain.MessagePagination) dto.MessagePagination {
	return dto.MessagePagination{
		Messages: c.MessagesDomainToDTO(pagination.Messages),
		Cursor:   pagination.Cursor,
	}
}

func (c chatConverter) ChatDomainToDTO(chat domain.Chat) dto.Chat {
	return dto.Chat{
		ID:              chat.ID,
		PhotoUrl:        getStringPointer(chat.PhotoUrl),
		FirstName:       chat.FirstName,
		LastName:        chat.LastName,
		LastMessage:     chat.LastMessage,
		TimeLastMessage: chat.TimeLastMessage,
	}
}

func (c chatConverter) ChatsDomainToDTO(chats []domain.Chat) []dto.Chat {
	result := make([]dto.Chat, len(chats))

	for i, chat := range chats {
		result[i] = c.ChatDomainToDTO(chat)
	}

	return result
}

func (c chatConverter) MessageGetToMessageCreate(message domain.MessageGet, isTrainer bool, userID int) domain.MessageCreate {
	var messageCreate domain.MessageCreate

	messageCreate.Message = getNullString(message.Message)
	messageCreate.ServiceID = getNullInt(message.ServiceID)
	messageCreate.IsToUser = isTrainer

	if isTrainer {
		messageCreate.UserID = message.To
		messageCreate.TrainerID = userID
	} else {
		messageCreate.UserID = userID
		messageCreate.TrainerID = message.To
	}

	fmt.Println(messageCreate)

	return messageCreate
}

func (c chatConverter) MessageCreateToMessage(message domain.MessageCreate, id int, t time.Time) domain.Message {
	return domain.Message{
		ID:        id,
		UserID:    message.UserID,
		TrainerID: message.TrainerID,
		Message:   getStringPointer(message.Message),
		ServiceID: getIntPointer(message.ServiceID),
		IsToUser:  message.IsToUser,
		Time:      t,
	}
}
