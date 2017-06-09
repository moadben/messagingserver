package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

//SocketServe is the function ran when ever a new request by a new client is created
func SocketServe(s *Server, w http.ResponseWriter, r *http.Request) {
	// Decoding the query parameter accompanying the request for websocket(holds username)
	values := r.URL.Query()
	// Creating a websocket between the user and server
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	newUser := User{Socket: conn, Username: values["username"][0]}
	s.Users[newUser.Username] = newUser
}

func main() {
	// Creating server object
	server := CreateServer()
	// Handler for new websocket request
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		SocketServe(server, w, r)
	})
	// Starting of the server
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error running server: " + err.Error())
	}
}
