package main

import (
	"errors"

	"github.com/gorilla/websocket"
)

// Server is the server object that handles Users
type Server struct {
	Users map[string]User
}

// Description stores a description that accompanies an offer/answer
type Description struct {
	User      string
	Recipient string
	Type      string
}

// Message holds a message that is read from a websocket
type Message struct {
	Type      string
	Recipient string
	Payload   string
	SDP       Description
}

// Connection stores the local and remote descriptions of a connection
// between two peers
type Connection struct {
	LocalDescription  Description
	RemoteDescription Description
}

// User holds the username and websocket for a new client
type User struct {
	Username    string
	Socket      *websocket.Conn
	Server      *Server
	Connections map[string]Connection
}

// Helper function to check if user exists and return it
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

// Helper function to check for legitimate connection between two peers
func checkConn(peerOne *User, peerTwo *User) bool {
	emptyDescription := Description{}
	peerOneConn := peerOne.Connections[peerTwo.Username]
	peerTwoConn := peerTwo.Connections[peerOne.Username]
	if peerOneConn.LocalDescription == peerTwoConn.RemoteDescription &&
		peerOneConn.RemoteDescription == peerTwoConn.LocalDescription &&
		peerOneConn.RemoteDescription != emptyDescription &&
		peerOneConn.LocalDescription != emptyDescription {
		return true
	}
	return false
}

// SetLocalDescription sets the local description of a connection between
// two peers
func (u *User) SetLocalDescription(m Message) Description {
	localDescription := Description{
		User:      u.Username,
		Recipient: m.Recipient,
		Type:      m.Type,
	}
	remoteDescription := u.Connections[m.Recipient].RemoteDescription
	u.Connections[m.Recipient] = Connection{LocalDescription: localDescription, RemoteDescription: remoteDescription}
	return localDescription
}

// SetRemoteDescription sets the remote description of a connection between
// two peers
func (u *User) SetRemoteDescription(m Message) {
	localDescription := u.Connections[m.SDP.User].LocalDescription
	u.Connections[m.SDP.User] = Connection{LocalDescription: localDescription, RemoteDescription: m.SDP}
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
	// Set Local Description
	description := u.SetLocalDescription(m)
	// Create an offer request
	newOffer = Message{"offerRequest", val.Username, u.Username + " Wants to create a connection", description}
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

	// Set Local Description
	description := u.SetLocalDescription(m)
	// Create an answer response
	newAnswer = Message{"answerResponse", val.Username, "accept offer", description}
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

	// See if there is a legitimate connection between the two peers
	if checkConn(u, &val) {
		// Create the message
		newMessage = Message{"message", val.Username, m.Payload, Description{}}
		return newMessage, nil
	}
	// Otherwise return an error
	return newMessage, errors.New("No connection between peers")
}

// CreateServer creates a new server object
func CreateServer() *Server {
	return &Server{
		Users: make(map[string]User),
	}
}
