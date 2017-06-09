package main

import (
	"errors"

	"github.com/gorilla/websocket"
)

// Server is the server object that handles Users
type Server struct {
	Users map[string]User
}

// Message holds a message that is read from a websocket
type Message struct {
	Type      string
	Recipient string
	Payload   string
}

// User holds the username and websocket for a new client
type User struct {
	Username string
	Socket   *websocket.Conn
	Server   *Server
}

func grabUser(u *User, m Message) (User, error) {
	var val User
	var ok bool
	var err error

	// Check to see if user exists
	if val, ok = u.Server.Users[m.Recipient]; !ok {
		err = errors.New("User not found")
	}
	return val, err
}

// CreateOffer takes the message and creates an Offer for the
// respective user
func (u *User) CreateOffer(m Message) (Message, error) {
	var newOffer Message

	// Check to see if user exists
	val, err := grabUser(u, m)
	if err != nil {
		return newOffer, errors.New("User not found")
	}

	// Create an offer request
	newOffer = Message{"offerRequest", val.Username, u.Username + " Wants to create a connection"}
	return newOffer, nil
}

// CreateAnswer takes the message and creates an Answer for the
// respective user
func (u *User) CreateAnswer(m Message) (Message, error) {
	var newAnswer Message

	// Check to see if user exists
	val, err := grabUser(u, m)
	if err != nil {
		return newAnswer, errors.New("User not found")
	}

	// Create an answer response
	newAnswer = Message{"answerResponse", val.Username, u.Username + " Accepts your connection"}
	return newAnswer, nil
}

// SendMessage takes the message and creates a message for the
// respective user
func (u *User) SendMessage(m Message) (Message, error) {
	var newMessage Message

	// Check to see if user exists
	val, err := grabUser(u, m)
	if err != nil {
		return newMessage, errors.New("User not found")
	}

	// Create the message
	newMessage = Message{"message", val.Username, m.Payload}
	return newMessage, nil
}

// CreateServer creates a new server object
func CreateServer() *Server {
	return &Server{
		Users: make(map[string]User),
	}
}
