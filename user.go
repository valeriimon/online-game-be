package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

type UserDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	conns     Set
}

func createUser(name string) *User {
	u := &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     "email@ex.com",
		Role:      "customer",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	app.addUser(u)

	return u
}

func (u *User) createConn(socket *websocket.Conn, root bool) *UserConnection {
	conn := &UserConnection{
		id:     uuid.New().String(),
		user:   u,
		socket: socket,
		root:   root,
		send:   make(chan *Message),
	}

	u.conns.add(conn)
	return conn
}

func (u *User) getRootConn() *UserConnection {
	rootConn, _ := u.conns.find(func(item T, index int) bool {
		return item.(UserConnection).root
	})

	return rootConn.(*UserConnection)
}

// @Deprecated
func (u *User) getDto() *UserDto {
	return &UserDto{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (u *User) getUserChats() []Chat {
	chats := make([]Chat, 0)

	app.chats.forEach(func(item T, i int) {
		for _, userID := range item.(Chat).Users {
			if userID == u.ID {
				chats = append(chats, item.(Chat))
				break
			}
		}
	})

	return chats
}

func (u *User) changeUserStatus(status string) {
	u.Status = status
}
