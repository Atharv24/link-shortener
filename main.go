package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"link-shortener/models"
	"link-shortener/services"
	"link-shortener/utils"
)

const (
	port = ":8000"
)

func main() {
	// Connect to database
	err := services.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer services.DisconnectDB()

	// Set up ticker to call DeleteExpiredLinks every hour
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			services.DeleteExpiredLinks()
		}
	}()

	// Set up HTTP handlers
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/login", services.OAuthLoginHandler)
	http.HandleFunc("/oauth2callback", services.OAuth2CallbackHandler)
	http.HandleFunc("/", redirectHandler)

	// Start HTTP server
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if !services.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check for valid OAuth 2.0 token
	if !services.IsAuthenticated(r) {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "Missing 'url' field", http.StatusBadRequest)
		return
	}

	parsedUrl, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Error parsing url: "+err.Error(), http.StatusBadRequest)
		return
	}

	// checking if scheme is not present and adding
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "http"
	}

	// Create new short URL
	shortURL := utils.GenerateShortURL()
	session, _ := utils.Store.Get(r, "session-name")
	email, _ := session.Values["email"].(string)
	l := models.Link{
		ID:         primitive.NewObjectID(),
		ShortURL:   shortURL,
		LongURL:    parsedUrl.String(),
		AddedBy:    email,
		CreatedAt:  time.Now(),
		AccessedAt: time.Now(),
	}

	// Insert new link into database
	err = services.AddLink(&l)
	if err != nil {
		http.Error(w, "Error inserting link into database", http.StatusInternalServerError)
		return
	}

	// Return new short URL
	json.NewEncoder(w).Encode(l)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	if shortURL == "favicon.ico" || shortURL == "" {
		return
	}
	// Check for valid OAuth 2.0 token
	if !services.IsAuthenticated(r) {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	l, err := services.GetLinkByShortURL(shortURL)
	if err != nil {
		http.Error(w, "Error retrieving link from database", http.StatusInternalServerError)
		return
	} else if l == nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Update the link's accessed at timestamp
	err = services.UpdateLinkAccessedAt(l)
	if err != nil {
		log.Print("Error updating links accessed at timestamp")
	}

	// Redirect to long URL
	http.Redirect(w, r, l.LongURL, http.StatusTemporaryRedirect)
}
