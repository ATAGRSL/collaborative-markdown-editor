package user

import (
	"fmt"
	"math/rand"
	"time"
)

// User represents a user in the collaborative editor
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Color     string    `json:"color"`
	CursorPos int       `json:"cursorPos"`
	LastSeen  time.Time `json:"lastSeen"`
}

// UserManager manages users across all rooms
type UserManager struct {
	users map[string]*User // userID -> User
	rooms map[string]map[string]*User // roomID -> userID -> User
}

// NewUserManager creates a new user manager
func NewUserManager() *UserManager {
	return &UserManager{
		users: make(map[string]*User),
		rooms: make(map[string]map[string]*User),
	}
}

// CreateUser creates a new user with a random color
func (um *UserManager) CreateUser(userID, username string) *User {
	// Generate random color
	colors := []string{
		"#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FFEAA7",
		"#DDA0DD", "#98D8C8", "#F7DC6F", "#BB8FCE", "#F8C471",
		"#85C1E9", "#F9E79F", "#A9DFBF", "#F5B7B1", "#C8B2DB",
	}
	color := colors[rand.Intn(len(colors))]

	user := &User{
		ID:        userID,
		Username:  username,
		Color:     color,
		CursorPos: 0,
		LastSeen:  time.Now(),
	}

	um.users[userID] = user
	return user
}

// GetUser gets a user by ID
func (um *UserManager) GetUser(userID string) (*User, bool) {
	user, exists := um.users[userID]
	return user, exists
}

// AddUserToRoom adds a user to a room
func (um *UserManager) AddUserToRoom(userID, roomID string) {
	if um.rooms[roomID] == nil {
		um.rooms[roomID] = make(map[string]*User)
	}
	if user, exists := um.users[userID]; exists {
		um.rooms[roomID][userID] = user
		user.LastSeen = time.Now()
	}
}

// RemoveUserFromRoom removes a user from a room
func (um *UserManager) RemoveUserFromRoom(userID, roomID string) {
	if roomUsers, exists := um.rooms[roomID]; exists {
		delete(roomUsers, userID)
	}
}

// GetRoomUsers gets all users in a room
func (um *UserManager) GetRoomUsers(roomID string) map[string]*User {
	roomUsers := make(map[string]*User)
	if users, exists := um.rooms[roomID]; exists {
		for userID, user := range users {
			// Create a copy to avoid race conditions
			roomUsers[userID] = &User{
				ID:        user.ID,
				Username:  user.Username,
				Color:     user.Color,
				CursorPos: user.CursorPos,
				LastSeen:  user.LastSeen,
			}
		}
	}
	return roomUsers
}

// UpdateUserCursor updates a user's cursor position
func (um *UserManager) UpdateUserCursor(userID string, position int) {
	if user, exists := um.users[userID]; exists {
		user.CursorPos = position
		user.LastSeen = time.Now()
	}
}

// RemoveUser removes a user completely
func (um *UserManager) RemoveUser(userID string) {
	delete(um.users, userID)
	for roomID := range um.rooms {
		delete(um.rooms[roomID], userID)
	}
}

// GenerateUsername generates a random username
func GenerateUsername() string {
	adjectives := []string{"Quick", "Lazy", "Happy", "Sad", "Fast", "Slow", "Bright", "Dark", "Cool", "Warm", "Big", "Small", "Tall", "Short", "Rich", "Poor"}
	nouns := []string{"Panda", "Tiger", "Eagle", "Shark", "Wolf", "Bear", "Lion", "Fox", "Owl", "Cat", "Dog", "Bird", "Fish", "Tree", "Star", "Moon"}

	return fmt.Sprintf("%s%s%d",
		adjectives[rand.Intn(len(adjectives))],
		nouns[rand.Intn(len(nouns))],
		rand.Intn(1000)+1)
}
