package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/jasbirnetwork/us/chat-app/hub"
	"github.com/jasbirnetwork/us/chat-app/models"

	"github.com/jasbirnetwork/us/chat-app/rabbitmq"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WSHandler upgrades the HTTP connection to a WebSocket and handles messaging.
func WSHandler(hubInstance *hub.Hub, getUsername func(r *http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := getUsername(r)
		if username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket Upgrade error: %v", err)
			return
		}
		client := &hub.Client{
			Conn:     conn,
			Username: username,
		}
		hubInstance.Register <- client

		defer func() {
			hubInstance.Unregister <- client
		}()

		// Listen for incoming messages.
		for {
			_, msgData, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}
			messageText := string(msgData)
			if strings.HasPrefix(messageText, "/stock=") {
				stockCode := strings.TrimSpace(strings.TrimPrefix(messageText, "/stock="))
				if stockCode != "" {
					if err := rabbitmq.PublishStockCommand(stockCode); err != nil {
						log.Printf("Error publishing stock command: %v", err)
					}
				}
				continue // Do not broadcast the command itself.
			}
			// Otherwise treat as a normal chat message.
			chatMsg := hub.CreateChatMessage(username, messageText)
			hub.AddMessage(chatMsg)
			wsMsg := models.WSMessage{
				Type: "chat",
				Data: chatMsg,
			}
			hubInstance.Broadcast <- wsMsg
		}
	}
}
