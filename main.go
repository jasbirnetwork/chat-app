package main

import (
	"log"
	"net/http"

	"github.com/jasbirnetwork/us/chat-app/handlers"
	"github.com/jasbirnetwork/us/chat-app/hub"
	"github.com/jasbirnetwork/us/chat-app/rabbitmq"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize RabbitMQ.
	if err := rabbitmq.Init(); err != nil {
		log.Fatalf("RabbitMQ initialization error: %v", err)
	}
	defer rabbitmq.Close()

	// Create and run the Hub.
	hubInstance := hub.NewHub()
	go hubInstance.Run()

	// Start consuming bot responses.
	if err := rabbitmq.ConsumeBotResponses(hubInstance.Broadcast, hub.AddMessage); err != nil {
		log.Fatalf("Failed to consume bot responses: %v", err)
	}

	// Setup HTTP routes.
	r := mux.NewRouter()
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler)
	r.HandleFunc("/chat", handlers.ChatHandler)
	r.HandleFunc("/ws", handlers.WSHandler(hubInstance, handlers.GetSessionUsername))
	r.HandleFunc("/search", handlers.SearchHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
	})

	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("HTTP Server error: %v", err)
	}
}
