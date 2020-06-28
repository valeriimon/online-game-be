package main

import "errors"

type Chat struct {
	ID        string   `json:"id"`
	ChatType  string   `json:"type"` // private | group
	Users     []string `json:"users"`
	CreatedBy string   `json:"created_by"`
	ChatLogo  string   `json:"chat_logo"`
}

type ChatUpdatedData struct {
	Chat     Chat          `json:"chat"`
	Messages []ChatMessage `json:"messages"`
}

func (chat *Chat) onMessage(m *ChatMessage, emittedBy *UserConnection) {
	app.addChatMessage(m)
	chat.onChatUpdated(emittedBy)
}

func (chat *Chat) onChatUpdated(emittedBy *UserConnection) {
	chatUpdatedData := ChatUpdatedData{
		Chat:     *chat,
		Messages: app.getChatMessagesByChatId(chat.ID),
	}
	msg, err := createMessage(chat.getComposedID(), chatUpdatedEvent, chatUpdatedData, emittedBy.user)
	if err != nil {
		app.logError("onChatMessage", err)
	}

	for _, userID := range chat.Users {
		u, err := app.getUserById(userID)
		if err != nil {
			app.logError("onChatMessage", err)
		}

		if u.Status == userStatusOffline {
			return
		}

		conn := u.getRootConn()
		if conn == nil {
			panic(errors.New("Root connection is missing when user is online"))
		}

		conn.send <- msg
	}
}

func (chat *Chat) getComposedID() string {
	return "chat: " + chat.ID
}
