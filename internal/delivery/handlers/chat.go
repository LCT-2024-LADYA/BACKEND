package handlers

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/delivery/middleware"
	"BACKEND/internal/services"
	"BACKEND/pkg/responses"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChatHandler struct {
	service   services.Chat
	converter converters.ChatConverter
}

func InitChatHandler(
	service services.Chat,
) *ChatHandler {
	return &ChatHandler{
		service:   service,
		converter: converters.InitChatConverter(),
	}
}

// GetUserChats
// @Summary Get User Chats
// @Description Get all chats for a user
// @Tags Chats
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param search query string false "Search term"
// @Success 200 {object} []dto.Chat "List of chats"
// @Failure 400 {object} responses.MessageResponse "Bad JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/chat/user [get]
func (h *ChatHandler) GetUserChats(c *gin.Context) {
	search := c.Query("search")
	userID := c.GetInt(middleware.UserID)

	chats, err := h.service.GetUserChats(c.Request.Context(), userID, search)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, chats)
}

// GetTrainerChats
// @Summary Get Trainer Chats
// @Description Get all chats for a trainer
// @Tags Chats
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param search query string false "Search term"
// @Success 200 {array} dto.Chat "List of chats"
// @Failure 400 {object} responses.MessageResponse "Bad JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/chat/trainer [get]
func (h *ChatHandler) GetTrainerChats(c *gin.Context) {
	search := c.Query("search")
	trainerID := c.GetInt(middleware.UserID)

	chats, err := h.service.GetTrainerChats(c.Request.Context(), trainerID, search)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, chats)
}

// GetChatMessageUser
// @Summary Get Chat Messages User
// @Description Get messages for a chat between a user and a trainer - user
// @Tags Chats
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param trainer_id path int true "Trainer ID"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.MessagePagination "List of messages with pagination"
// @Failure 400 {object} responses.MessageResponse "Bad path or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/chat/user/{trainer_id} [get]
func (h *ChatHandler) GetChatMessageUser(c *gin.Context) {
	trainerID, err := strconv.Atoi(c.Param("trainer_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	cursorStr := c.Query("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	userID := c.GetInt(middleware.UserID)

	messages, err := h.service.GetChatMessage(c.Request.Context(), userID, trainerID, cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// GetChatMessageTrainer
// @Summary Get Chat Messages Trainer
// @Description Get messages for a chat between a user and a trainer - trainer
// @Tags Chats
// @Accept json
// @Produce json
// @Param access_token header string true "Access token"
// @Param user_id path int true "User ID"
// @Param cursor query int false "Cursor for pagination"
// @Success 200 {object} dto.MessagePagination "List of messages with pagination"
// @Failure 400 {object} responses.MessageResponse "Bad path or JWT provided"
// @Failure 401 {object} responses.MessageResponse "JWT is expired or invalid"
// @Failure 500 "Internal server error"
// @Router /api/chat/trainer/{user_id} [get]
func (h *ChatHandler) GetChatMessageTrainer(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadPath})
		return
	}

	cursorStr := c.Query("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.MessageResponse{Message: responses.ResponseBadQuery})
		return
	}

	trainerID := c.GetInt(middleware.UserID)

	messages, err := h.service.GetChatMessage(c.Request.Context(), userID, trainerID, cursor)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, messages)
}
