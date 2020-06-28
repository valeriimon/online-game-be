package main

import "encoding/json"

type Message struct {
	Address   string          `json:"address"`
	Type      string          `json:"type"`
	Body      json.RawMessage `json:"body"`
	CreatedBy *User           `json:"createdBy"`
	extraData T
	conn      *UserConnection
}

type ErrorData struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

type RequestNewGameData struct {
	GameType       string   `json:"gameType"`
	InvitedUsersID []string `json:"invitedUsersIds"`
	WelcomeMessage string   `json:"welcomeMessage"`
}

type InvitationData struct {
	InviteToID     string `json:"inviteTo"`
	WelcomeMessage string `json:"welcomeMessage"`
	PlaceInfo      string `json:"place"`
}

type UserChangeData struct {
	userID string `json:"user_id"`
}

type HubData struct {
	ID         string `json:"id"`
	UserOnline int    `json:"online"`
	Users      []User `json:"usersList"`
}

type LobbyData struct {
	ID         string `json:"id"`
	Users      []User `json:"users"`
	GameType   string
	GameConfig T `json:"game_config"`
}

type AppData struct {
	Hubs    []Hub   `json:"hubs"`
	Lobbies []Lobby `json:"lobbies"`
	Users   []User  `json:"users"`
}

type GlobalChatMessage struct {
	Text string `json:"message_text"`
}

type PrivateMessage struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Text      string `json:"message_text"`
	CreatedAt int    `json:"created_at"`
	Read      bool   `json:"read"`
}

type GroupChatMessage struct {
	From      string `json:"from"`
	Text      string `json:"message_text"`
	createdAt int    `json: created_at`
	Read      bool   `json:"read"`
}

type CreateGroupChatData struct {
	Users []string `json:"users"`
}

type GroupChatData struct {
	Users     []User    `json:"users"`
	Messages  []Message `json:"messages"`
	CreatedBy User      `json:"created_by"`
}

type ChatMessage struct {
	ID          string `json:"id"`
	ChatID      string `json:"chat_id"`
	MessageType string `json:"message_type"` // text | link | image
	From        string `json:"from"`
	Message     string `json:"message"`
}

func createMessage(adr, eventType string, body T, createdBy *User) (*Message, error) {
	msgBody, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	return &Message{
		Address:   adr,
		Body:      json.RawMessage(msgBody),
		CreatedBy: createdBy,
		Type:      eventType,
	}, nil
}

func createErrorData(code int, reason string) *ErrorData {
	return &ErrorData{
		code,
		reason,
	}
}

func createHubData(id string, online int, Users []User) *HubData {
	return &HubData{
		ID:         id,
		UserOnline: online,
		Users:      Users,
	}
}

func createLobbyData(id string, Users []User, gameType string, gameConfig T) *LobbyData {
	return &LobbyData{
		ID:         id,
		Users:      Users,
		GameType:   gameType,
		GameConfig: gameConfig,
	}
}

func createInvitationData(inviteTo string, welcomeMessage string) *InvitationData {
	return &InvitationData{
		InviteToID:     inviteTo,
		WelcomeMessage: welcomeMessage,
		PlaceInfo:      "coming soon",
	}
}

func (msg *Message) parseBody() T {
	var msgBody T
	err := json.Unmarshal(msg.Body, &msgBody)
	if err != nil {
		sendError("server", 1, "Parsing body has failed", msg.conn)
		return nil
	}

	return msgBody
}

func sendError(adr string, code int, reason string, sendTo *UserConnection) {
	errBody := ErrorData{code, reason}
	sendMessage(adr, errorMessageEvent, errBody, nil, sendTo)
}

func sendMessage(adr, eventType string, body T, createdBy *User, sendTo *UserConnection) {
	msg, err := createMessage(adr, eventType, body, createdBy)
	if err != nil {
		app.logError("sendMessage", err)
		return
	}

	sendTo.send <- msg
}
