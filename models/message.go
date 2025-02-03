package models

import "time"

// ChatMessage represents a single chat message.
type ChatMessage struct {
	ID        int       `json:"id"`
	User      string    `json:"user"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// StockRequest is the JSON payload for a stock command.
type StockRequest struct {
	StockCode string `json:"stock_code"`
}

// BotResponse is the JSON payload returned by the bot.
type BotResponse struct {
	Message string `json:"message"`
}

// WSMessage is the structure used for WebSocket communication.
type WSMessage struct {
	Type string      `json:"type"` // "chat" or "userlist"
	Data interface{} `json:"data"`
}
