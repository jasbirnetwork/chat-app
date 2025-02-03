package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	amqp091_go "github.com/rabbitmq/amqp091-go"
)

type StockRequest struct {
	StockCode string `json:"stock_code"`
}

type BotResponse struct {
	Message string `json:"message"`
}

const (
	stockRequestQueue = "stockRequests"
	botResponseQueue  = "botResponses"
	rabbitURL         = "amqp://guest:guest@localhost:5672/"
)

// fetchStockQuote calls the external API and parses the CSV response.
func fetchStockQuote(stockCode string) (string, error) {
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", strings.ToLower(stockCode))
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error calling stock API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response: %v", err)
	}

	reader := csv.NewReader(bytes.NewReader(body))
	records, err := reader.ReadAll()
	if err != nil {
		return "", fmt.Errorf("error parsing CSV: %v", err)
	}

	if len(records) < 2 {
		return "", fmt.Errorf("unexpected CSV format")
	}

	headers := records[0]
	values := records[1]
	var symbolIdx, closeIdx int = -1, -1
	for i, h := range headers {
		switch strings.ToLower(h) {
		case "symbol":
			symbolIdx = i
		case "close":
			closeIdx = i
		}
	}
	if symbolIdx == -1 || closeIdx == -1 {
		return "", fmt.Errorf("required columns not found")
	}
	quote := values[closeIdx]
	if quote == "N/D" || quote == "" {
		return "", fmt.Errorf("invalid stock code")
	}
	if _, err := strconv.ParseFloat(quote, 64); err != nil {
		return "", fmt.Errorf("invalid price value")
	}
	stockSymbol := strings.ToUpper(values[symbolIdx])
	message := fmt.Sprintf("%s quote is $%s per share", stockSymbol, quote)
	return message, nil
}

func main() {
	conn, err := amqp091_go.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Bot: failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Bot: failed to open channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		stockRequestQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Bot: failed to declare stockRequests queue: %v", err)
	}
	_, err = ch.QueueDeclare(
		botResponseQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Bot: failed to declare botResponses queue: %v", err)
	}

	msgs, err := ch.Consume(
		stockRequestQueue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Bot: failed to register consumer: %v", err)
	}

	log.Println("Bot: Waiting for stock requests...")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var req StockRequest
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Bot: error decoding message: %v", err)
				d.Ack(false)
				continue
			}
			log.Printf("Bot: Received stock request for '%s'", req.StockCode)
			responseMsg, err := fetchStockQuote(req.StockCode)
			if err != nil {
				responseMsg = fmt.Sprintf("Could not fetch stock for %s", req.StockCode)
			}
			botResp := BotResponse{Message: responseMsg}
			body, err := json.Marshal(botResp)
			if err != nil {
				log.Printf("Bot: error marshaling bot response: %v", err)
				d.Ack(false)
				continue
			}
			err = ch.Publish(
				"",
				botResponseQueue,
				false,
				false,
				amqp091_go.Publishing{
					ContentType: "application/json",
					Body:        body,
				},
			)
			if err != nil {
				log.Printf("Bot: error publishing bot response: %v", err)
			} else {
				log.Printf("Bot: Published response: %s", responseMsg)
			}
			d.Ack(false)
		}
	}()

	<-forever
}
