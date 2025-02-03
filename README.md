# Chat App

A real‐time chat application in **Go** that showcases:

- **RabbitMQ** integration for triggering stock quotes via a separate bot worker.
- **WebSocket** connections to broadcast messages instantly.
- **In‐memory** chat history (last 50 messages).
- A **WhatsApp‐style** UI with a fixed left sidebar and scrollable chat panel.

## Table of Contents

1. [Features](#features)
2. [Screenshots](#screenshots)
3. [Architecture Overview](#architecture-overview)
4. [Dependencies](#dependencies)
5. [Project Structure](#project-structure)
6. [Installation & Setup](#installation--setup)
7. [Running the Application](#running-the-application)
8. [Using the Chat](#using-the-chat)
9. [Searching Messages](#searching-messages)
10. [Stock Bot](#stock-bot)
11. [License](#license)

---

## Features

- **User Login**: Simple login with a chosen username (no password checks).
- **Real-Time Chat**: Uses WebSockets to push new messages to all clients instantly.
- **RabbitMQ Bot**: `/stock=...` commands are published to a queue. A bot worker fetches quotes from [Stooq](https://stooq.com) and posts back.
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
