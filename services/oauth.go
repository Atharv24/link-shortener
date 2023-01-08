package services

import (
	"os"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
)

func InitOAuth(credentialsFile string) error {
	// Load OAuth 2.0 credentials from file
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return err
	}

	// Parse OAuth 2.0 credentials
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		return err
	}
	oauthConfig = config

	return nil
}

func IsAuthenticated(r *http.Request) bool {
	// // Get Auth token
	// token := r.Header.Get("Authorization")

	// // Set up HTTP client
	// client := oauthConfig.Client(context.Background())

	// // Send request to OAuth 2.0 server to verify token
	// resp, err := client.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token)
	// if err != nil {
	// 	log.Printf("Error verifying token: %v", err)
	// 	return false
	// }
	// defer resp.Body.Close()

	// // Check response status
	// if resp.StatusCode != http.StatusOK {
	// 	log.Printf("Token verification failed: %d", resp.StatusCode)
	// 	return false
	// }

	return true
}