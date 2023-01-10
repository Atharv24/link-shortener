package services

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"link-shortener/models"
	"link-shortener/utils"
)

var (
	oauthConfig *oauth2.Config
)

func init() {
	// Load OAuth 2.0 credentials from file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Error reading credentials file")
	}

	// Parse OAuth 2.0 credentials
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		log.Fatalf("Error parsing credentials")
	}
	oauthConfig = config
}

func IsAuthenticated(r *http.Request) bool {
	// Get access token from request header
	accessToken := r.Header.Get("Authorization")

	if(accessToken == "") {
		// Get access token from session
		session, _ := utils.Store.Get(r, "session-name")
		accessToken, ok := session.Values["access_token"].(string)
		if !ok || accessToken == "" {
			return false
		}
	}

	//Set up http client
	client := oauthConfig.Client(context.Background(), &oauth2.Token{AccessToken: accessToken})

	// Send request to OAuth 2.0 server to verify token
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/tokeninfo")
	if err != nil {
		log.Printf("Error verifying token: %v", err)
		return false
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Token verification failed: %d", resp.StatusCode)
		return false
	}

	return true
}

func OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authorization URL
	authURL := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)

	// Redirect the user to the authorization URL
	http.Redirect(w, r, authURL, http.StatusFound)
}

func OAuth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Get the authorization code from the query parameters
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	// Create an HTTP client with the access token
	client := oauthConfig.Client(context.Background(), token)

	// Use the HTTP client to make a request to the Google API
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "failed to get user info", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "failed to parse response body", http.StatusInternalServerError)
		return
	}

	session, _ := utils.Store.Get(r, "session-name")
	session.Values["access_token"] = token.AccessToken
	session.Values["email"] = user.Email
	session.Save(r, w)
	log.Printf("saved token: " + token.AccessToken)

	err = AddUser(&user)
	if(err != nil) {
		log.Println("Error adding user to DB. Email: " + user.Email)
	}
	// Redirect to home
	http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
}
