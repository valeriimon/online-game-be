package main

import "fmt"

type App struct {
	lobbies      Set // []*Lobby
	hubs         Set // []*Hub
	users        Set // []*Users
	chats        Set
	chatMessages Set
	news         Set
}

type Entity interface {
	addUserConn(User *UserConnection)
	removeUserConn(User *UserConnection)
	onMessage(m *Message)
}

func (app *App) onMessage(m *Message) {
	switch m.Type {
	case getNewsEvent:
		{
			news := getNews()
			sendMessage("server", getNewsEvent, news, nil, m.conn)
			break
		}
	case getUserChatsEvent:
		{
			u := m.CreatedBy
			chats := u.getUserChats()
			sendMessage("server", getUserChatsEvent, chats, nil, m.conn)
			break
		}
	case getUserEvent:
		{
			// code
			break
		}
	case getUserFriendsEvent:
		{
			sendMessage("server", getUserFriendsEvent, app.users, nil, m.conn)
		}
	}
}

// @User repo
func (app *App) addUser(c *User) {
	app.users.add(c)
}

// @User repo
func (app *App) getUserById(id string) (*User, error) {
	c, err := app.users.find(func(c T, i int) bool {
		return c.(*User).ID == id
	})

	if err != nil {
		return nil, err
	}

	return c.(*User), nil
}

// @User repo
func (app *App) removeUserById(id string) bool {
	c, err := app.users.find(func(c T, i int) bool {
		return c.(*User).ID == id
	})

	if err != nil {
		return false
	}

	return app.users.delete(c)
}

func (app *App) addHub(h *Hub) {
	app.hubs.add(h)
	app.sendAppInfoMsg()
}

func (app *App) addLobby(l *Lobby) {
	app.lobbies.add(l)
	app.sendAppInfoMsg()
}

// @Chat repo
func (app *App) addChat(chat *Chat) {
	app.chats.add(chat)
}

// @Chat repo
func (app *App) getChatByID(chatID string) (*Chat, error) {
	chat, err := app.chats.find(func(item T, i int) bool {
		return item.(Chat).ID == chatID
	})

	if err != nil {
		return nil, err
	}

	return chat.(*Chat), nil
}

// @Chat repo
func (app *App) removeChatByID(chatID string) bool {
	chat, err := app.chats.find(func(item T, i int) bool {
		return item.(Chat).ID == chatID
	})

	if err != nil {
		return false
	}

	app.removeChatMessagesByChatId(chatID)
	return app.chats.delete(chat)
}

// @ChatMessages repo
func (app *App) addChatMessage(m *ChatMessage) {
	app.chats.add(m)
}

// @ChatMessages repo
func (app *App) getChatMessagesByChatId(chatID string) []ChatMessage {
	msgs := make([]ChatMessage, 0)

	app.chatMessages.forEach(func(item T, i int) {
		if item.(ChatMessage).ChatID == chatID {
			msgs = append(msgs, item.(ChatMessage))
		}
	})

	return msgs
}

// @ChatMessages repo
func (app *App) removeChatMessagesByChatId(chatID string) {
	msgs := app.getChatMessagesByChatId(chatID)

	for _, msg := range msgs {
		app.chatMessages.delete(msg)
	}
}

// @ChatMessages repo
func (app *App) removeMessageByID(msgID string) bool {
	msg, err := app.chatMessages.find(func(item T, i int) bool {
		return item.(ChatMessage).ID == msgID
	})

	if err != nil {
		return false
	}

	return app.chatMessages.delete(msg)
}

func (app *App) sendAppInfoMsg() {
	msg := app.getAppInfoMsg()
	app.broadcastToSpecificRoles("admin", msg, "server_updated")
}

func (app *App) getAppInfoMsg() *AppData {
	hubs := make([]Hub, 0)
	lobbies := make([]Lobby, 0)
	users := make([]User, 0)

	app.hubs.forEach(func(h T, i int) {
		hubs = append(hubs, *h.(*Hub))
	})

	app.lobbies.forEach(func(l T, i int) {
		lobbies = append(lobbies, *l.(*Lobby))
	})

	app.users.forEach(func(c T, i int) {
		users = append(users, *c.(*User))
	})

	return &AppData{
		hubs,
		lobbies,
		users,
	}
}

func (app *App) broadcastToSpecificRoles(role, body T, eventType string) {
	app.users.forEach(func(item T, index int) {
		c := item.(User)
		if c.Role == role {
			conn := c.getRootConn()
			app.sendUserMessage(conn, body, eventType)
		}
	})
}

func (app *App) sendUserMessage(conn *UserConnection, body T, eventType string) {
	msg, err := createMessage("server", eventType, body, conn.getUser())

	if err != nil {
		app.logError("app: sendUserMessage", err)
	}

	conn.send <- msg
}

func (app *App) logError(method string, err error) {
	fmt.Printf("INTERNAL [ "+method+" ] error - %v", err.Error())
}
