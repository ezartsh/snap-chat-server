package websockets

import (
	"encoding/json"
	"slices"
)

type Sender struct {
	Name     string
	Username string
}

type ClientMessage struct {
	RoomUID     string
	RoomName    string
	IsError     bool
	MessageType string
	Target      []string
	Sender      Sender
	Message     []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Inbound message from for specific client
	chatIn chan ClientMessage
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		chatIn:     make(chan ClientMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.chatIn:
			for client, ok := range h.clients {
				if ok {
					if slices.Contains(message.Target, client.auth.Username) {
						byteMessage, _ := json.Marshal(message)
						select {
						case client.send <- byteMessage:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
