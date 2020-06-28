package main

import "errors"

type eventHandler func(payload *Message, listeners []*UserConnection)

type WS struct {
	events         map[string][]eventHandler
	eventListeners map[string][]*UserConnection
	guard          func(conn *UserConnection) bool
	getConns       func() []*UserConnection
}

// Add handler for specified event message
func (ws *WS) on(eventType string, handler eventHandler) {
	ws.events[eventType] = append(ws.events[eventType], handler)
}

// Trigger emitted event handlers
func (ws *WS) handleEvent(message *Message) {
	handlers := ws.findEventHandlers(Url(message.Type))

	if handlers {
		return
	}

	listeners, _ := ws.eventListeners[message.Type]

	for _, h := range handlers {
		h(message, listeners)
	}
}

// On upcoming message handler
func (ws *WS) onMessage(message *Message) {
	ok := ws.guard(message.conn)
	if ok {
		ws.handleEvent(message)
	}
}

// Broadcast message to all registered Users
func (ws *WS) broadcast(message *Message) {
	conns := ws.getConns()
	for _, conn := range conns {
		conn.send <- message
	}
}

// Adds event listener to event
func (ws *WS) addEventListener(eventType string, conn *UserConnection) error {
	_, ok := ws.events[eventType]
	if !ok {
		return errors.New("Event is not registered")
	}

	ws.eventListeners[eventType] = append(ws.eventListeners[eventType], conn)

	return nil
}

func (ws *WS) findEventHandlers(eventType Url) []eventHandler {
	for event, handlers := range ws.events {
		u := Url(event)
		if u.isMatched(eventType) {
			return handlers
		}
	}

	return nil
}
