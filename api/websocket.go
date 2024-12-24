package api

import (
	"fmt"
	"github/Shubhpreet-Rana/projects/internal/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type WebSocketManager struct {
	hub *Hub
}

func NewWebSocketHandler() *WebSocketManager {
	hub := NewHub()
	go hub.Run() // Start the Hub's main loop
	return &WebSocketManager{hub: hub}
}

// NewWebSockets initializes WebSocket routes and connects them to the Hub.
func (h *WebSocketManager) RegisterSockets(app *fiber.App) {
	// Middleware to upgrade requests to WebSocket
	app.Use("/ws", h.WebSocketUpgradeMiddleware())
	// WebSocket route handler
	app.Get("/ws", websocket.New(h.WebSocketHandler))
}

// WebSocketUpgradeMiddleware checks if the request is a WebSocket upgrade.
func (h *WebSocketManager) WebSocketUpgradeMiddleware() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

// WebSocketHandler handles WebSocket connection and message processing.
func (h *WebSocketManager) WebSocketHandler(c *websocket.Conn) {
	logging.InfoLogger.Println("WebSocket connection established")
	defer func() {
		logging.InfoLogger.Println("WebSocket connection closed")
		c.Close()
	}()

	roomName := c.Query("room_id", "defaultRoom")
	userName := c.Query("name", "User")

	// Create and register the client in the Hub
	client := &Client{
		hub:  h.hub,
		conn: c,
		send: make(chan []byte, 256),
		room: roomName,
	}

	h.hub.register <- client

	// Join a room based on ROOMID (for now, it's the default room)
	room := roomName // Default room for the client
	logging.InfoLogger.Printf("Client joined room: %s", room)
	h.hub.joinRoom(client, room)

	// Auto-reply to all members in the room with a welcome message
	h.hub.broadcast <- BroadcastRequest{Message{
		Room:    room,
		Content: []byte(fmt.Sprintf("Welcome %s to the room! %s", userName, room)),
	}, client}

	// Start reading messages from the WebSocket
	go client.writePump()
	client.readPump()
}
