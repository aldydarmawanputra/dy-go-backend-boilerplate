package realtime

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Upgrade rejects non-WebSocket requests before the handler runs.
func Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

// Handler returns the Fiber WebSocket handler bound to this hub.
func (h *Hub) Handler() fiber.Handler {
	return websocket.New(func(conn *websocket.Conn) {
		client := &Client{hub: h, conn: conn, send: make(chan []byte, 16)}
		h.register <- client

		go client.writePump()
		client.readPump()
	})
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
		// Incoming messages are ignored in this general skeleton.
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
