package hub

import (
	"encoding/json"
	"log"
	"time"

	"collaborative-markdown-editor/internal/client"
	"collaborative-markdown-editor/internal/ot"
	"collaborative-markdown-editor/internal/user"
)

// RegisterRequest represents a client registration request
type RegisterRequest struct {
	Client *client.Client
	RoomID string
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients in each room.
type Hub struct {
	// Registered clients organized by room ID
	rooms map[string]map[*client.Client]bool

	// OT managers for each room
	otManagers map[string]*ot.Manager

	// User manager
	userManager *user.UserManager

	// Inbound messages from the clients
	broadcast chan client.Message

	// Register requests from the clients
	register chan RegisterRequest

	// Unregister requests from clients
	unregister chan *client.Client
}

// NewHub creates a new hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:    make(chan client.Message),
		register:     make(chan RegisterRequest),
		unregister:   make(chan *client.Client),
		rooms:        make(map[string]map[*client.Client]bool),
		otManagers:   make(map[string]*ot.Manager),
		userManager:  user.NewUserManager(),
	}
}

// Run starts the hub and handles client registration, unregistration, and message broadcasting
func (h *Hub) Run() {
	for {
		select {
		case request := <-h.register:
			// Initialize room if it doesn't exist
			if h.rooms[request.RoomID] == nil {
				h.rooms[request.RoomID] = make(map[*client.Client]bool)
				// Create OT manager for new room
				h.otManagers[request.RoomID] = ot.NewManager("")
			}
			// Add client to the room
			h.rooms[request.RoomID][request.Client] = true
			log.Printf("Client registered in room %s. Total clients in room: %d", request.RoomID, len(h.rooms[request.RoomID]))

			// Broadcast updated user list to all clients in the room
			h.broadcastUserList(request.RoomID)

		case client := <-h.unregister:
			// Remove client from all rooms
			for roomID, clients := range h.rooms {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					log.Printf("Client unregistered from room %s. Remaining clients: %d", roomID, len(clients))

					// Broadcast updated user list to all remaining clients in the room
					h.broadcastUserList(roomID)

					// Clean up empty rooms
					if len(clients) == 0 {
						delete(h.rooms, roomID)
						log.Printf("Room %s deleted (empty)", roomID)
					}
					break
				}
			}

		case message := <-h.broadcast:
			// Process the message with OT manager
			if otManager, ok := h.otManagers[message.RoomID]; ok {
				msgContent := string(message.Content)

				// Check if this is a JSON operation or plain content
				if len(msgContent) > 0 && msgContent[0] == "{" {
					// JSON operation
					operation, err := ot.OperationFromJSON(message.Content)
					if err != nil {
						log.Printf("Failed to parse operation: %v", err)
						continue
					}

					// Apply the operation
					transformedOp, err := otManager.ApplyOperation(operation)
					if err != nil {
						log.Printf("Failed to apply operation: %v", err)
						continue
					}

					if transformedOp == nil {
						// Operation was cancelled due to conflicts
						continue
					}

					// Broadcast the transformed operation to all clients in the room except sender
					if clients, ok := h.rooms[message.RoomID]; ok {
						transformedJSON, _ := transformedOp.ToJSON()
						for client := range clients {
							if client.ID != message.ClientID {
								select {
								case client.Send <- transformedJSON:
								default:
									// Client's send channel is full or closed, remove client
									close(client.Send)
									delete(clients, client)
								}
							}
						}
					}

					// Update cursor position for the user who made the change
					if user, exists := h.userManager.GetUser(message.ClientID); exists {
						user.CursorPos = operation.Position + len(operation.Character)
						user.LastSeen = time.Now()
					}
				} else {
					// Plain text content
					otManager.SetCurrentDocument(msgContent)

					// Broadcast the plain content to all clients in the room except sender
					if clients, ok := h.rooms[message.RoomID]; ok {
						for client := range clients {
							if client.ID != message.ClientID {
								select {
								case client.Send <- message.Content:
								default:
									// Client's send channel is full or closed, remove client
									close(client.Send)
									delete(clients, client)
								}
							}
						}
					}

					// Update cursor position for the user who made the change
					if user, exists := h.userManager.GetUser(message.ClientID); exists {
						user.CursorPos = len(msgContent)
						user.LastSeen = time.Now()
					}
				}
			}
		}
	}
}

// GetRoomContent returns the current content of a room from OT manager
func (h *Hub) GetRoomContent(roomID string) string {
	if otManager, ok := h.otManagers[roomID]; ok {
		return otManager.GetCurrentDocument()
	}
	return ""
}

// GetUserManager returns the user manager
func (h *Hub) GetUserManager() *user.UserManager {
	return h.userManager
}

// GetRoomUsers returns all users in a room
func (h *Hub) GetRoomUsers(roomID string) map[string]*user.User {
	return h.userManager.GetRoomUsers(roomID)
}

// Register adds a client to the hub
func (h *Hub) Register(client *client.Client) {
	request := RegisterRequest{
		Client: client,
		RoomID: client.RoomID,
	}
	h.register <- request
}

// Broadcast sends a message to all clients in a room except the sender
func (h *Hub) Broadcast(msg client.Message) {
	h.broadcast <- msg
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *client.Client) {
	h.unregister <- client
}

// broadcastUserList sends the current user list to all clients in a room
func (h *Hub) broadcastUserList(roomID string) {
	if clients, ok := h.rooms[roomID]; ok {
		roomUsers := h.userManager.GetRoomUsers(roomID)

		// Convert users to a simple format for JSON
		type UserInfo struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Color    string `json:"color"`
		}

		var userList []UserInfo
		for _, user := range roomUsers {
			userList = append(userList, UserInfo{
				ID:       user.ID,
				Username: user.Username,
				Color:    user.Color,
			})
		}

		// Create a special message type for user list updates
		type UserListMessage struct {
			Type string     `json:"type"`
			Users []UserInfo `json:"users"`
		}

		userListMsg := UserListMessage{
			Type: "userList",
			Users: userList,
		}

		jsonData, err := json.Marshal(userListMsg)
		if err != nil {
			log.Printf("Failed to marshal user list: %v", err)
			return
		}

		// Send to all clients in the room
		for client := range clients {
			select {
			case client.Send <- jsonData:
			default:
				// Client's send channel is full or closed, remove client
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}
