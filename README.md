# Chat App

A real‐time chat application in **Go** that showcases:

- **RabbitMQ** integration for triggering stock quotes via a separate bot worker.
- **WebSocket** connections to broadcast messages instantly.
- **In‐memory** chat history (last 50 messages).
- A **WhatsApp‐style** UI with a fixed left sidebar and scrollable chat panel.

## Table of Contents

- [Chat App](#chat-app)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Screenshots](#screenshots)
  - [Architecture Overview](#architecture-overview)
  - [Dependencies](#dependencies)
  - [Project Structure](#project-structure)
  - [Installation \& Setup](#installation--setup)
  - [Running the Application](#running-the-application)
  - [Using the Chat](#using-the-chat)
  - [Stock Bot](#stock-bot)
  - [Future / TODO Features](#future--todo-features)

---

## Features

- **User Login**: Simple login with a chosen username (no password checks).
- **Real-Time Chat**: Uses WebSockets to push new messages to all clients instantly.
- **RabbitMQ Bot**: `/stock=XYZ` commands are published to a queue. A bot worker fetches quotes from [Stooq](https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv) and posts back.
- **Fixed Layout**: Left sidebar (connected users, search box) and right chat panel, each scrolling internally only when content grows large.

---

## Screenshots

1. **Login Page**  
   ![Login Page](docs/images/login.png)

2. **Main Chat UI**  
   ![Chat UI](docs/images/chat_ui.png)

3. **Stock Quote Example**  
   ![Stock Quote](docs/images/stock_quote.png)

---

## Architecture Overview

[Browser 1] <--(WebSocket)--> [Go Chat Server] <--(AMQP)--> [RabbitMQ] <--(AMQP)--> [Bot Worker]
[Browser 2] <-----------------------------------^

1. **Go Chat Server**

   - Serves HTML, CSS, JS to the client.
   - Maintains a list of connected WebSocket clients.
   - Publishes stock requests to RabbitMQ.

2. **Bot Worker**

   - Subscribes to `stockRequests` queue for `/stock=...` commands.
   - Fetches quotes from Stooq, publishes a result message to `botResponses`.

3. **RabbitMQ**
   - Routes messages between the chat server and the bot worker.

---

## Dependencies

- **Go 1.20+**
- **RabbitMQ** (listening on `amqp://guest:guest@localhost:5672`)
- **Libraries** listed in `go.mod`:
  - [gorilla/mux](https://github.com/gorilla/mux)
  - [gorilla/websocket](https://github.com/gorilla/websocket)
  - [gorilla/sessions](https://github.com/gorilla/sessions)
  - [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go)

---

## Project Structure

```bash
    chat-app/
    ├── bot/
    │ └── bot.go # Bot Worker
    ├── handlers/
    │ ├── login.go
    │ └── ws.go
    ├── hub/
    │ └── hub.go
    ├── models/
    │ └── message.go
    ├── rabbitmq/
    │ └── rabbitmq.go
    ├── static/
    │ ├── css/
    │ │ └── styles.css
    │ └── js/
    │ └── app.js
    ├── templates/
    │ ├── login.html
    │ └── chat.html
    ├── tests/
    │ └── stock_test.go
    ├── main.go
    └── go.mod
```

---

## Installation & Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/jasbirnetwork/chat-app.git
   cd chat-app

   ```

2. **Install Go Modules**

   ```bash
   go mod tidy

   ```

3. **Start RabbitMQ**
   - Ensure RabbitMQ is running locally on `(amqp://guest:guest@localhost:5672)`.
   - Install from [RabbitMQ official site](https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv) or use any method you prefer.

---

## Running the Application

1. **Start the Chat Server**

   ```bash
   go run main.go
   ```

   - The server listens on [http://localhost:8080](http://localhost:8080)

2. **Start the Bot Worker (in a new terminal)**

   ```bash
    cd bot
    go run bot/bot.go

   ```

3. **Open the Chat**
   - Go to [http://localhost:8080](http://localhost:8080/login).
   - Enter a username and click Enter Chat.

---

## Using the Chat

1. **Multiple Users**

   - Open multiple browser tabs or incognito windows. Log in with different usernames.

2. **Send a Message**
   Type your message in the input box, press Enter or Send.

3. **Stock Command**
   ```bash
   /stock=AAPL.US
   ```
   - After a short delay, a bot message shows the current stock price.

---

## Stock Bot

1. **Queues**
   - `stockRequests`: Chat server publishes a JSON payload with the `stock_code`.
   - `botResponses`: The bot publishes its results, which the server consumes and broadcasts as a message from `bot`.

- CSV API: The bot calls Stooq.com to fetch the quote.

---

## Future / TODO Features

1. **Queues**
   - **User Authentication**: Implement real password checks and user management.
   - **Private Messaging**: Let users open direct 1‐on‐1 chats.
   - **Multiple Chat Rooms**: Expand beyond a single “global” room.
   - **Persistent Storage**: Use a database (Postgres, MySQL) instead of an in‐memory history.
   - **Typing Indicators**: Show when a user is typing.
   - **Message Reactions**: Add emojis or reactions to existing messages
