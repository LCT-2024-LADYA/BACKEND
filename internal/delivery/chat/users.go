package chat

import (
	"BACKEND/internal/models/domain"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
)

type IncomeMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type User struct {
	id        int
	isTrainer bool
	conn      *websocket.Conn
	server    *Server
	income    chan *domain.Message
	done      chan bool
}

func NewUser(id int, isTrainer bool, conn *websocket.Conn, server *Server) *User {
	return &User{
		id:        id,
		isTrainer: isTrainer,
		conn:      conn,
		server:    server,
		income:    make(chan *domain.Message, 100),
		done:      make(chan bool, 100),
	}
}

func (u *User) GetStarted() {
	go u.listen()
	go u.write()
}

func (u *User) listen() {
	for {
		select {
		case message := <-u.income:
			err := u.conn.WriteJSON(message)
			if err != nil {
				u.server.err(err)
				u.Done()
				return
			}
		case <-u.done:
			u.server.delUser(u)
			return
		}
	}
}

func (u *User) write() {
	for {
		select {
		case <-u.done:
			u.server.delUser(u)
			return
		default:
			var incomeMessage IncomeMessage
			err := u.conn.ReadJSON(&incomeMessage)
			if err != nil {
				u.server.err(err)
				u.Done()
				return
			} else {
				switch incomeMessage.Type {
				case "message":
					messageBytes, err := json.Marshal(incomeMessage.Data)
					if err != nil {
						u.server.err(err)
						u.Done()
						return
					}

					var messageGet domain.MessageGet
					if err = json.Unmarshal(messageBytes, &messageGet); err != nil {
						u.server.err(err)
						u.Done()
						return
					}

					// Создание записи в БД
					messageCreate := u.server.converter.MessageGetToMessageCreate(messageGet, u.isTrainer, u.id)

					createdID, time, err := u.server.service.CreateMessage(context.Background(), messageCreate)
					if err != nil {
						u.server.err(err)
						u.Done()
						return
					}

					message := u.server.converter.MessageCreateToMessage(messageCreate, createdID, time)

					u.server.sendMessage(&message)
				}
			}
		}
	}
}

func (u *User) Done() {
	u.done <- true
}
