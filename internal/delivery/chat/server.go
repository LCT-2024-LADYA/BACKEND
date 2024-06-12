package chat

import (
	"BACKEND/internal/converters"
	"BACKEND/internal/models/domain"
	"BACKEND/internal/services"
	"BACKEND/pkg/responses"
	"BACKEND/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	mu             sync.Mutex
	users          []*User
	messagesToSend chan *domain.Message
	addUsers       chan *User
	delUsers       chan *User
	errs           chan error

	service   services.Chat
	converter converters.ChatConverter
	jwtUtil   utils.JWT
	logger    zerolog.Logger
}

func NewServer(
	service services.Chat,
	jwtUtil utils.JWT,
	logger zerolog.Logger,
) *Server {
	return &Server{
		mu:             sync.Mutex{},
		users:          []*User{},
		messagesToSend: make(chan *domain.Message, 100),
		addUsers:       make(chan *User, 100),
		delUsers:       make(chan *User, 100),
		errs:           make(chan error, 100),
		service:        service,
		converter:      converters.InitChatConverter(),
		jwtUtil:        jwtUtil,
		logger:         logger,
	}
}

func (s *Server) addUser(user *User) {
	s.addUsers <- user
}

func (s *Server) delUser(user *User) {
	s.delUsers <- user
}

func (s *Server) sendMessage(message *domain.Message) {
	s.messagesToSend <- message
}

func (s *Server) err(err error) {
	s.errs <- err
}

func (s *Server) send(message *domain.Message) {
	var userTo *User

	if message.IsToUser {
		for _, user := range s.users {
			if !user.isTrainer && user.id == message.UserID {
				userTo = user
				break
			}
		}
	} else {
		for _, user := range s.users {
			if user.isTrainer && user.id == message.TrainerID {
				userTo = user
				break
			}
		}
	}

	if userTo == nil {
		if message.IsToUser {
			s.logger.Info().Msg(fmt.Sprintf("User %d is offline", message.UserID))
		} else {
			s.logger.Info().Msg(fmt.Sprintf("User %d is offline", message.TrainerID))
		}
		return
	}

	userTo.income <- message
	if message.IsToUser {
		s.logger.Info().Msg(fmt.Sprintf("Message to user %d is send", message.UserID))
	} else {
		s.logger.Info().Msg(fmt.Sprintf("Message to user %d is send", message.TrainerID))
	}
}

func (s *Server) Listen() {
	s.logger.Info().Msg("Start listen ...")
	for {
		select {
		case user := <-s.addUsers:
			s.mu.Lock()
			s.users = append(s.users, user)
			s.mu.Unlock()
			s.logger.Info().Msg(fmt.Sprintf("User %d (is trainer: %t) added to server\n", user.id, user.isTrainer))
		case user := <-s.delUsers:
			for i, delUser := range s.users {
				if delUser.id == user.id {
					s.mu.Lock()
					s.users = append(s.users[:i], s.users[i+1:]...)
					s.mu.Unlock()
					delUser.conn.Close()
					s.logger.Info().Msg(fmt.Sprintf("User %d (is trainer: %t) deleted from server\n", user.id, user.isTrainer))
					break
				}
			}
		case err := <-s.errs:
			s.logger.Error().Msg(fmt.Sprintf("WS error: %s\n", err.Error()))
		case message := <-s.messagesToSend:
			s.send(message)
		}
	}
}

func (s *Server) ChatHandler(c *gin.Context) {
	accessToken := c.GetHeader("Sec-WebSocket-Protocol")
	if accessToken == "" {
		s.logger.Error().Msg("No jwt provided")
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.MessageResponse{Message: "No jwt provided"})
		return
	}

	userData, isValid, err := s.jwtUtil.Authorize(accessToken, utils.User, utils.Trainer)
	if err != nil {
		s.logger.Error().Msg(fmt.Sprintf("Troubles while getting user info from jwt: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, responses.MessageResponse{Message: "Bad jwt provided"})
		return
	}

	if !isValid {
		s.logger.Error().Msg("Access token is expired or invalid or user type mismatch")
		c.AbortWithStatusJSON(http.StatusUnauthorized, responses.MessageResponse{Message: "Access token is expired or invalid"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, http.Header{
		"Sec-WebSocket-Protocol": []string{accessToken},
	})

	if err != nil {
		s.errs <- err
	}

	s.logger.Info().Msg(fmt.Sprintf("User %d connected", userData.ID))

	var isTrainer bool
	if userData.UserType == utils.Trainer {
		isTrainer = true
	}

	user := NewUser(userData.ID, isTrainer, conn, s)

	s.addUser(user)
	user.GetStarted()
}
