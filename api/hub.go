package api

import (
	"github/Shubhpreet-Rana/projects/internal/logging"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool            // Registered clients
	rooms      map[string]map[*Client]bool // Room-to-clients mapping
	broadcast  chan BroadcastRequest       // Broadcast channel for messages
	register   chan *Client                // Channel for registering new clients
	unregister chan *Client                // Channel for unregistering clients
	mu         sync.Mutex                  // Protects the clients and rooms
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan BroadcastRequest),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
}

func (h *Hub) joinRoom(client *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.rooms[room]; !ok {
		h.rooms[room] = make(map[*Client]bool)
	}
	h.rooms[room][client] = true
}

func (h *Hub) leaveRoom(client *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if clients, ok := h.rooms[room]; ok {
		delete(clients, client)
	}
}

func (h *Hub) broadcastMsg(message Message, sender *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if the room exists
	if clients, ok := h.rooms[message.Room]; ok {
		for client := range clients {
			// Skip the sender
			if client == sender {
				continue
			}
			select {
			case client.send <- message.Content:
				// Successfully sent the message to the client
			default:
				// If the client's send buffer is full, close the channel and remove the client
				close(client.send)
				delete(clients, client)
				logging.ErrorLogger.Printf("Send buffer full for client %s. Disconnecting.", client.conn.RemoteAddr())
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
			// Optionally, clean up room memberships here
			for room, clients := range h.rooms {
				if _, ok := clients[client]; ok {
					h.leaveRoom(client, room)
				}
			}
		case broadcastRequest := <-h.broadcast:
			h.broadcastMsg(broadcastRequest.Message, broadcastRequest.Sender)
		}
	}
}
