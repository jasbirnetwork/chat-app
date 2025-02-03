package hub

import (
	"log"
	"sync"
	"time"

	"github.com/jasbirnetwork/us/chat-app/models"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client.
type Client struct {
	Conn     *websocket.Conn
	Username string
}

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan models.WSMessage
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// NewHub creates and returns a new Hub.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan models.WSMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run listens for register/unregister events and broadcasts messages.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			h.broadcastUserList()
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.mu.Unlock()
			h.broadcastUserList()
		case message := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Printf("WebSocket send error: %v", err)
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// broadcastUserList sends an updated list of connected usernames to all clients.
func (h *Hub) broadcastUserList() {
	h.mu.Lock()
	defer h.mu.Unlock()
	var users []string
	for client := range h.Clients {
		users = append(users, client.Username)
	}
	msg := models.WSMessage{
		Type: "userlist",
		Data: users,
	}
	for client := range h.Clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Printf("User list send error: %v", err)
			client.Conn.Close()
			delete(h.Clients, client)
		}
	}
}

var (
	nextMessageID int = 1
	chatHistory       = make([]models.ChatMessage, 0)
	chatHistoryMu     sync.Mutex
)

// AddMessage adds a new chat message to the history (keeping only the last 50).
func AddMessage(msg models.ChatMessage) {
	chatHistoryMu.Lock()
	defer chatHistoryMu.Unlock()
	msg.ID = nextMessageID
	nextMessageID++
	chatHistory = append(chatHistory, msg)
	if len(chatHistory) > 50 {
		chatHistory = chatHistory[len(chatHistory)-50:]
	}
}

// GetHistory returns the current chat history.
func GetHistory() []models.ChatMessage {
	chatHistoryMu.Lock()
	defer chatHistoryMu.Unlock()
	return chatHistory
}

// CreateChatMessage is a helper to create a chat message.
func CreateChatMessage(username, content string) models.ChatMessage {
	return models.ChatMessage{
		User:      username,
		Content:   content,
		Timestamp: time.Now(),
	}
}
