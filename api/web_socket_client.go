package api

import (
	"github/Shubhpreet-Rana/projects/internal/logging"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	room string
}

// ReadPump reads messages from the WebSocket connection and forwards them to the hub
func (c *Client) readPump() {
	defer func() {
		logging.InfoLogger.Println("Closing connection for client:", c.conn.RemoteAddr())
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			logging.ErrorLogger.Printf("Error reading message: %v", err)
			break
		}

		// Log the received message
		logging.InfoLogger.Printf("Received message from %s: %s", c.conn.RemoteAddr(), string(message))

		room := c.room // For simplicity, assume all messages are for the default room
		c.hub.broadcast <- BroadcastRequest{Message: Message{
			Room:    room,
			Content: message,
		}, Sender: c}

		// // Optionally, send an auto-reply to the sender
		// autoReply := fmt.Sprintf("Received your message: %s", message)
		// c.send <- []byte(autoReply)
	}
}

// WritePump writes messages to the WebSocket connection from the send channel
func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			logging.ErrorLogger.Println("Write error:", err)
			break
		}
	}

}
