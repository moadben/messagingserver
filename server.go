package main

import "github.com/gorilla/websocket"

// User holds the username and websocket for a new client
type User struct {
	Username string `json:"username"`
	Socket   *websocket.Conn
}

// Server is the server object that handles Users
type Server struct {
	Users map[string]User
}

// CreateServer creates a new server object
func CreateServer() *Server {
	return &Server{
		Users: make(map[string]User),
	}
}
