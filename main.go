package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var utils = Utils{}

var app App

func main() {
	fmt.Println("Hello")

	app = App{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	hub := createHub(true)

	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/complete-user-auth", func(w http.ResponseWriter, r *http.Request) {
		// code
	})

	http.HandleFunc("/signup", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// code
		r.ParseMultipartForm(0)
		user := &User{
			ID:        uuid.New().String(),
			Name:      r.FormValue("name"),
			Email:     r.FormValue("email"),
			Role:      r.FormValue("role"),
			Password:  r.FormValue("password"),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		app.addUser(user)
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}))

	http.HandleFunc("/login", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(0)
		user, err := app.users.find(func(item T, i int) bool {
			return item.(*User).Password == r.FormValue("password") && item.(*User).Email == r.FormValue("email")
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user.(*User))
	}))

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// code
	})

	http.HandleFunc("/ws-root", func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err.Error())
		}

		matches, ok := r.URL.Query()["userId"]
		if !ok {
			socket.SetCloseHandler(func(code int, text string) error {
				return errors.New("Cannot find user")
			})
			socket.Close()
			return
		}

		user, err := app.getUserById(matches[0])
		if err != nil {
			socket.SetCloseHandler(func(code int, text string) error {
				return errors.New("Cannot find user")
			})
			socket.Close()
			return
		}
		conn := user.createConn(socket, true)

		go conn.read()
		go conn.write()
		hub.addUserConn(conn)
	})

	// http.HandleFunc("ws/:room_id", func(w http.ResponseWriter, r *http.Request) {
	// 	conn, err := upgrader.Upgrade(w, r, nil)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}

	// 	hub := createHub()

	// 	client := &Client{
	// 		Name: "Valerii",
	// 		conn: conn,
	// 		send: make(chan *Message),
	// 	}

	// 	// conn.

	// 	go client.read()
	// 	go client.write()

	// 	hub.register <- client
	// })
	http.ListenAndServe(":8080", nil)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
