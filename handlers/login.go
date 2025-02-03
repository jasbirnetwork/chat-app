package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/jasbirnetwork/us/chat-app/hub"
	"github.com/jasbirnetwork/us/chat-app/models"

	"github.com/gorilla/sessions"
)

var sessionStore = sessions.NewCookieStore([]byte("super-secret-key"))
var templates = template.Must(template.ParseGlob("templates/*.html"))

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    searchTerm := strings.TrimSpace(r.URL.Query().Get("q"))
    log.Printf("Search term: %v", searchTerm)

    // If empty, return an empty slice (not nil)
    if searchTerm == "" {
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode([]models.ChatMessage{})
        return
    }

    allMessages := hub.GetHistory() // e.g. []ChatMessage
    log.Printf("All messages: %+v", allMessages)

    var filtered []models.ChatMessage
    lowerTerm := strings.ToLower(searchTerm)
    for _, msg := range allMessages {
        if strings.Contains(strings.ToLower(msg.Content), lowerTerm) {
            filtered = append(filtered, msg)
        }
    }

    log.Printf("Filtered messages: %+v", filtered)

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(filtered)
}

// GetSessionUsername retrieves the username from session.
func GetSessionUsername(r *http.Request) string {
	session, err := sessionStore.Get(r, "session")
	if err != nil {
		return ""
	}
	username, ok := session.Values["username"].(string)
	if !ok {
		return ""
	}
	return username
}

// LoginHandler handles GET and POST requests for /login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		if strings.TrimSpace(username) == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		session, _ := sessionStore.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// LogoutHandler clears the session and redirects to /login.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "session")
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ChatHandler renders the chat page along with the current chat history.
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	username := GetSessionUsername(r)
	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	data := struct {
		Username string
		History  []models.ChatMessage
	}{
		Username: username,
		History:  hub.GetHistory(),
	}
	templates.ExecuteTemplate(w, "chat.html", data)
}
