package main

import "github.com/google/uuid"

type Room struct {
	id    string
	conns Set
	ws    *WS
}

func createRoom() *Room {
	r := &Room{
		id:    uuid.New().String(),
		conns: Set{},
	}

	ws := &WS{
		events:   make(map[string][]eventHandler),
		getConns: r.getConns,
		guard:    func(conn *UserConnection) bool { return true },
	}

	r.ws = ws

	return r
}

func (r *Room) getConns() []*UserConnection {
	conns := make([]*UserConnection, 0)
	r.conns.forEach(func(c T, i int) {
		conns = append(conns, c.(*UserConnection))
	})

	return conns
}

func (r *Room) broadcast(m *Message) {
	r.ws.broadcast(m)
}

func (r *Room) sendTo(conn *UserConnection) {

}
