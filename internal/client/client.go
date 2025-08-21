package client

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"collaborative-markdown-editor/internal/ot"
	"collaborative-markdown-editor/internal/user"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		return true
	},
}

// Client represents a WebSocket connection to a single client
type Client struct {
	// The WebSocket connection
	Conn *websocket.Conn

	// Buffered channel of outbound messages
	Send chan []byte

	// Current content of the document
	CurrentContent string

	// Previous content for diff calculation
	PreviousContent string

	// Room ID this client belongs to
	RoomID string

	// Unique client ID
	ID string

	// Current version for OT
	Version int

	// User information
	User *user.User
}

// Message represents a message from a client
type Message struct {
	RoomID   string `json:"roomId"`
	ClientID string `json:"clientId"`
	Content  []byte `json:"content"`
}



// NewClient creates a new client instance
func NewClient(conn *websocket.Conn, roomID, clientID string, user *user.User) *Client {
	return &Client{
		Conn:            conn,
		Send:            make(chan []byte, 256),
		CurrentContent:  "",
		PreviousContent: "",
		RoomID:          roomID,
		ID:              clientID,
		Version:         0,
		User:            user,
	}
}

// readPump pumps messages from the WebSocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump(hub interface {
	Broadcast(msg Message)
	Unregister(client *Client)
}) {
	defer func() {
		hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msgContent := string(message)

		// Check if this is a JSON operation or plain content
		if len(msgContent) > 0 && msgContent[0] == "{" {
			// Try to parse as JSON operation
			operation, err := ot.OperationFromJSON(message)
			if err != nil {
				log.Printf("Failed to parse operation: %v", err)
				return
			}

			// Update current content based on operation
			c.applyOperation(operation)
			c.Version = operation.Version

			// Send the operation to hub
			msg := Message{
				RoomID:   c.RoomID,
				ClientID: c.ID,
				Content:  message,
			}

			hub.Broadcast(msg)
		} else {
			// Plain text content (from JavaScript)
			newContent := msgContent

			// Simple approach: replace entire content
			c.PreviousContent = c.CurrentContent
			c.CurrentContent = newContent
			c.Version++

			// Send the full content to hub
			msg := Message{
				RoomID:   c.RoomID,
				ClientID: c.ID,
				Content:  []byte(newContent),
			}

			hub.Broadcast(msg)
		}
	}
}



// writePump pumps messages from the hub to the WebSocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Parse the operation and apply to local content
			operation, err := ot.OperationFromJSON(message)
			if err != nil {
				log.Printf("Failed to parse operation: %v", err)
				continue
			}

			// Apply operation to local content
			c.applyOperation(operation)

			// Send the full current content to the client
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(c.CurrentContent))

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// applyOperation applies an OT operation to the client's current content
func (c *Client) applyOperation(operation *ot.Operation) {
	runes := []rune(c.CurrentContent)

	switch operation.Type {
	case ot.Insert:
		if operation.Position > len(runes) {
			operation.Position = len(runes)
		}
		newRunes := make([]rune, 0, len(runes)+utf8.RuneCountInString(operation.Character))
		newRunes = append(newRunes, runes[:operation.Position]...)
		newRunes = append(newRunes, []rune(operation.Character)...)
		newRunes = append(newRunes, runes[operation.Position:]...)
		c.CurrentContent = string(newRunes)

	case ot.Delete:
		if operation.Position >= len(runes) {
			return
		}
		end := operation.Position + operation.Length
		if end > len(runes) {
			end = len(runes)
		}
		newRunes := make([]rune, 0, len(runes)-operation.Length)
		newRunes = append(newRunes, runes[:operation.Position]...)
		newRunes = append(newRunes, runes[end:]...)
		c.CurrentContent = string(newRunes)
	}
}
