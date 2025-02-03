package rabbitmq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jasbirnetwork/us/chat-app/models"

	amqp091_go "github.com/rabbitmq/amqp091-go"
)

const (
	StockRequestQueue = "stockRequests"
	BotResponseQueue  = "botResponses"
	RabbitURL         = "amqp://guest:guest@localhost:5672/"
)

var (
	conn    *amqp091_go.Connection
	channel *amqp091_go.Channel
)

// Init establishes a connection to RabbitMQ and declares queues.
func Init() error {
	var err error
	conn, err = amqp091_go.Dial(RabbitURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	_, err = channel.QueueDeclare(
		StockRequestQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare stockRequests queue: %v", err)
	}
	_, err = channel.QueueDeclare(
		BotResponseQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare botResponses queue: %v", err)
	}
	return nil
}

// PublishStockCommand publishes a stock command message to RabbitMQ.
func PublishStockCommand(stockCode string) error {
	req := models.StockRequest{StockCode: stockCode}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	err = channel.Publish(
		"", // default exchange
		StockRequestQueue,
		false,
		false,
		amqp091_go.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

// ConsumeBotResponses consumes bot responses and forwards them via the provided broadcast channel.
func ConsumeBotResponses(broadcast chan models.WSMessage, addMessageFunc func(models.ChatMessage)) error {
	msgs, err := channel.Consume(
		BotResponseQueue,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %v", err)
	}
	go func() {
		for d := range msgs {
			var resp models.BotResponse
			if err := json.Unmarshal(d.Body, &resp); err != nil {
				continue
			}
			chatMsg := models.ChatMessage{
				User:      "bot",
				Content:   resp.Message,
				Timestamp: time.Now(),
			}
			addMessageFunc(chatMsg)
			wsMsg := models.WSMessage{
				Type: "chat",
				Data: chatMsg,
			}
			broadcast <- wsMsg
		}
	}()
	return nil
}

// Close closes the RabbitMQ channel and connection.
func Close() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
