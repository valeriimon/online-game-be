package main

import "github.com/google/uuid"

import "fmt"

type Lobby struct {
	id          string
	ws          *WS
	events      map[string][]eventHandler
	conns       Set // []*UserConnection
	gameType    string
	gameConfigs T
}

func createLobby() *Lobby {
	l := &Lobby{
		id:       uuid.New().String(),
		gameType: "",
		events:   make(map[string][]eventHandler),
		conns:    Set{},
	}

	l.ws = &WS{
		events: make(map[string][]eventHandler),
		guard: func(conn *UserConnection) bool {
			return true
		},
		getConns: func() []*UserConnection {
			conns := make([]*UserConnection, 0)

			l.conns.forEach(func(c T, i int) {
				conns = append(conns, c.(*UserConnection))
			})

			return conns
		},
	}

	app.addLobby(l)
	l.registerEvents()

	fmt.Println("lobby created")

	return l
}

func (l *Lobby) registerEvents() {
	l.ws.on(updateLobbyEvent, l.updateLobby)
}

func (l *Lobby) addUserConn(conn *UserConnection) {
	l.conns.add(conn)
	conn.registerConn(l)
	l.onLobbyUpdated()
}

func (l *Lobby) onMessage(m *Message) {}

func (l *Lobby) removeUserConn(conn *UserConnection) {
	l.conns.delete(conn)
	l.onLobbyUpdated()
}

func (l *Lobby) updateLobby(message *Message) {
	body := message.parseBody()
	if body == nil {
		return
	}

	l.gameConfigs = body.(*LobbyData).GameConfig

	l.onLobbyUpdated()
}

func (l *Lobby) onLobbyUpdated() {
	lobbyInfoMsg := createLobbyData(l.id, l.getUsers(), l.gameType, l.gameConfigs)

	msg, err := createMessage("lobby: "+l.id, lobbyUpdatedEvent, lobbyInfoMsg, nil)

	if err != nil {
		panic(err.Error())
	}

	l.ws.broadcast(msg)
}

func (l *Lobby) getUsers() []User {
	Users := make([]User, 0)

	l.conns.forEach(func(c T, i int) {
		User := c.(*UserConnection).getUser()
		Users = append(Users, *User)
	})

	return Users
}
