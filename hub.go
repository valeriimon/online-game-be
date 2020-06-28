package main

import (
	"fmt"
	"github.com/google/uuid"
)

type Hub struct {
	id    string
	root  bool
	ws    *WS
	conns Set // []*UserConnection
	rooms Set

	message    chan *Message
	register   chan *UserConnection
	unregister chan *UserConnection
}

func createHub(isRoot bool) *Hub {
	hub := &Hub{
		id:         uuid.New().String(),
		root:       isRoot,
		conns:      Set{},
		message:    make(chan *Message),
		register:   make(chan *UserConnection),
		unregister: make(chan *UserConnection),
	}

	hub.ws = &WS{
		events:         make(map[string][]eventHandler),
		eventListeners: make(map[string][]*UserConnection),
		guard:          hub.messageGuard,
		getConns:       hub.getConns,
	}

	app.addHub(hub)
	hub.registerEvents()

	fmt.Println("hub created")

	return hub
}

func (hub *Hub) run() {
	for {
		select {
		case conn := <-hub.register:
			hub.registerConn(conn)

		case conn := <-hub.unregister:
			hub.unregisterUserConn(conn)

		case message := <-hub.message:
			hub.ws.onMessage(message)
		}
	}
}

// @Interface
func (hub *Hub) onMessage(m *Message) {
	hub.message <- m
}

// @Interface
func (hub *Hub) addUserConn(conn *UserConnection) {
	hub.register <- conn
}

// @Interface
func (hub *Hub) removeUserConn(conn *UserConnection) {
	hub.unregister <- conn
}

// Register all possible events and their handlers
func (hub *Hub) registerEvents() {
	hub.ws.on(userUpdatedEvent, hub.onUserUpdated)

	// Chat handlers
	hub.ws.on(globalChatMessageEvent, hub.onGlobalChatMessage)
	hub.ws.on(newChatEvent, hub.onNewChat)
	hub.ws.on(chatMessageEvent, hub.onChatMessage)
	hub.ws.on(chatUpdatedEvent, hub.onChatUpdated)

	// News handlers
	hub.ws.on(getNewsEvent, hub.onGetNews)
}

func (hub *Hub) registerConn(conn *UserConnection) {
	hub.conns.add(conn)
	conn.registerConn(hub)
	if hub.root {
		conn.user.changeUserStatus(userStatusOnline)
	}

	hub.onHubUpdated()
}

func (hub *Hub) unregisterUserConn(conn *UserConnection) {
	hub.conns.delete(conn)
	if hub.root {
		conn.user.changeUserStatus(userStatusOffline)
	}

	hub.onHubUpdated()
}

func (hub *Hub) messageGuard(conn *UserConnection) bool {
	if _, err := hub.conns.get(conn); err != nil {
		sendError(hub.getComposedID(), 1, "User is not registered in Hub", conn)
		return false
	}

	return true
}

func (hub *Hub) getConns() []*UserConnection {
	conns := make([]*UserConnection, 0)
	hub.conns.forEach(func(c T, i int) {
		conns = append(conns, c.(*UserConnection))
	})

	return conns
}

func (hub *Hub) onHubUpdated() {
	users := hub.getUsers()
	hubInfo := createHubData(hub.id, len(users), users)

	msg, err := createMessage(hub.getComposedID(), hubUpdatedEvent, hubInfo, nil)
	if err != nil {
		panic(err.Error())
	}
	hub.ws.broadcast(msg)
}

func (hub *Hub) onUserUpdated(m *Message) {
	hub.onHubUpdated()
}

// @Root
func (hub *Hub) onGetNews(m *Message) {
	news := getNews()

	sendMessage(hub.getComposedID(), getNewsEvent, news, m.conn.user, m.conn)
}

// @Root
func (hub *Hub) onGlobalChatMessage(m *Message) {
	hub.ws.broadcast(m)
}

func (hub *Hub) onNewChat(m *Message) {
	body := m.parseBody()
	if body == nil {
		return
	}

	chat := body.(*Chat)
	chat.ID = uuid.New().String()
	app.addChat(chat)
	sendMessage(hub.getComposedID(), newChatEvent, chat, m.conn.user, m.conn)
}

func (hub *Hub) onChatMessage(m *Message) {
	body := m.parseBody()
	if body == nil {
		return
	}

	chat, err := app.getChatByID(body.(ChatMessage).ChatID)
	if err != nil {
		sendError(hub.getComposedID(), 500, err.Error(), m.conn)
		return
	}

	chat.onMessage(body.(*ChatMessage), m.conn)
}

func (hub *Hub) onChatUpdated(m *Message) {
	body := m.parseBody()
	if body == nil {
		return
	}

	chat, err := app.getChatByID(body.(*ChatUpdatedData).Chat.ID)
	if err != nil {
		sendError(hub.getComposedID(), 500, err.Error(), m.conn)
		return
	}

	chat.onChatUpdated(m.conn)
}

// func (hub *Hub) onCreateLobby(message *Message) {
// 	lobby := createLobby()
// 	lobby.addUserConn(message.conn)

// 	lobbyInfo := createLobbyData(lobby.id, lobby.getUsers(), lobby.gameType, lobby.gameConfigs)

// 	hub.sendMessage(lobbyCreateEvent, nil, lobbyInfo, message.conn)
// }

func (hub *Hub) getUsers() []User {
	users := make([]User, 0)

	hub.conns.forEach(func(c T, i int) {
		users = append(users, *c.(*UserConnection).user)
	})

	return users
}

func (hub *Hub) getConnectionByUserID(userID string) *UserConnection {
	conn, err := hub.conns.find(func(item T, i int) bool {
		return item.(UserConnection).user.ID == userID
	})

	if err != nil {
		return nil
	}

	return conn.(*UserConnection)
}

func (hub *Hub) getComposedID() string {
	return "hub: " + hub.id
}
