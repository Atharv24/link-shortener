package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"link-shortener/models"
	"link-shortener/services"
	"link-shortener/utils"
)

var (
	shortURLCache *cache.Cache
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

	// Create short URL cache - caching the short urls for 1 day (if not accessed again)
	shortURLCache = cache.New(24*time.Hour, 10*time.Minute)

	// Set up HTTP handlers
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	// Start HTTP server
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
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

	// URL has not been shortened, create new short URL
	shortURL := utils.GenerateShortURL()
	l := models.Link{
		ID: primitive.NewObjectID(),
		ShortURL:  shortURL,
		LongURL:   req.URL,
		CreatedAt: time.Now(),
	}

	// Insert new link into database
	err := services.AddLink(&l)
	if err != nil {
		http.Error(w, "Error inserting link into database", http.StatusInternalServerError)
		return
	}

	// Return new short URL
	json.NewEncoder(w).Encode(l)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// Check for valid OAuth 2.0 token
	if !services.IsAuthenticated(r) {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	shortURL := r.URL.Path[1:]
	log.Println(shortURL)

	// Check cache for short URL
	val, found := shortURLCache.Get(shortURL)
	if found {
		// Redirect to long URL
		http.Redirect(w, r, val.(string), http.StatusTemporaryRedirect)
		return
	}

	// Short URL not in cache, check database
	l, err := services.GetLinkByShortURL(shortURL)
	if err != nil {
		http.Error(w, "Error retrieving link from database", http.StatusInternalServerError)
		return
	} else if l == nil {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	// Redirect to long URL
	http.Redirect(w, r, l.LongURL, http.StatusTemporaryRedirect)

	// Add short URL to cache
	shortURLCache.Set(l.ShortURL, l.LongURL, cache.DefaultExpiration)
}
