// main.go contains the main function that listens and serves for our server,
// and the functions that listen on a websocket.
package main

import (
	"fmt"
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
	// Create a new user object
	newUser := User{Socket: conn, Username: values["username"][0], Server: s, Connections: make(map[string]Connection)}
	s.Users[newUser.Username] = newUser

	Listen(&newUser)
}

// Listen listens on a websocket in a new goroutine and sends/reads messages.
func Listen(user *User) {
	// Defer closing of socket in case of emergency exit
	defer func() {
		delete(user.Server.Users, user.Username)
		user.Socket.Close()
	}()

	// infinite loop to read from opened socket
	for {
		// Read the message
		m := Message{}
		err := user.Socket.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
			continue
		}

		//Switch statement to check message type and react accordingly
		switch {
		// if user wants to make an offer
		case m.Type == "offer":
			offer, err := user.CreateOffer(m)
			if err != nil {
				fmt.Println("Error creating offer", err)
				continue
			}
			if err = user.Server.Users[offer.Recipient].Socket.WriteJSON(offer); err != nil {
				fmt.Println(err)
				continue
			} else {
				// Set Remote Description for recipient
				u := user.Server.Users[offer.Recipient]
				u.SetRemoteDescription(offer)
			}

		// if user wants to answer an offer
		case m.Type == "answer":
			answer, err := user.CreateAnswer(m)
			if err != nil {
				fmt.Println("Error creating answer", err)
				continue
			}
			if err = user.Server.Users[answer.Recipient].Socket.WriteJSON(answer); err != nil {
				fmt.Println(err)
				continue
			} else {
				// Set Remote Description for offer-er
				u := user.Server.Users[answer.Recipient]
				u.SetRemoteDescription(answer)
			}

		// if user wants to send a message
		case m.Type == "message":
			message, err := user.CreateMessage(m)
			if err != nil {
				fmt.Println("Error creating message", err)
				continue
			}
			if err = user.Server.Users[message.Recipient].Socket.WriteJSON(message); err != nil {
				fmt.Println(err)
				continue
			}

		// if users want to exchange ICE candiates
		case m.Type == "ice-candidate":
			message, err := user.CreateICECandidate(m)
			if err != nil {
				fmt.Println("Error creating message", err)
				continue
			}
			if err = user.Server.Users[message.Recipient].Socket.WriteJSON(message); err != nil {
				fmt.Println(err)
				continue
			}
		// if user wants to disconnect
		case m.Type == "logout":
			delete(user.Server.Users, user.Username)
			user.Socket.Close()
			return
		}
	}
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
