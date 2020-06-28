package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type UserConnection struct {
	id           string
	user         *User
	registeredAt Entity
	root         bool
	socket       *websocket.Conn
	send         chan *Message
}

// Redundant method
func (conn *UserConnection) getUser() *User {
	return conn.user
}

func (conn *UserConnection) registerConn(entity Entity) {
	conn.registeredAt = entity
}

func (conn *UserConnection) destroy() {
	conn.registeredAt.removeUserConn(conn)
}

func (conn *UserConnection) write() {
	defer func() {
		fmt.Println("closing write")
		conn.socket.Close()
	}()

	for {
		message := <-conn.send
		conn.writeMessage(message)
	}
}

func (conn *UserConnection) read() {
	defer func() {
		conn.destroy()
		conn.socket.Close()
	}()

	for {
		err := conn.readMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("disconnect")
			}

			break
		}
	}
}

func (conn *UserConnection) readMessage() error {
	var message Message

	err := conn.socket.ReadJSON(&message)
	if err != nil {
		return err
	}

	message.CreatedBy = conn.user
	message.conn = conn

	if message.Address == "server" {
		app.onMessage(&message)
	} else {
		conn.registeredAt.onMessage(&message)
	}

	return nil
}

func (conn *UserConnection) writeMessage(message *Message) error {
	err := conn.socket.WriteJSON(message)
	if err != nil {
		return err
	}

	return nil
}
