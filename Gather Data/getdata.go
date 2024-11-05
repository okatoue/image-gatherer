package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

var oauthConfig *oauth2.Config
var URL string

func set_oauth(input string) *oauth2.Config {
	if input == "google" {
		CLIENT_ID := "GOOGLE_CLIENT_ID"
		CLIENTSECRET := "GOOGLE_CLIENT_SECRET"
		URL = "https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=10"

		return &oauth2.Config{
			ClientID:     os.Getenv(CLIENT_ID),
			ClientSecret: os.Getenv(CLIENTSECRET),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
			RedirectURL:  "http://localhost:8080/callback",
		}
	} else {
		CLIENT_ID := "FACEBOOK_CLIENT_ID"
		CLIENTSECRET := "FACEBOOK_CLIENT_SECRET"
		URL = "https://graph.facebook.com/me/photos?type=uploaded&fields=id,name,images"

		return &oauth2.Config{
			ClientID:     os.Getenv(CLIENT_ID),
			ClientSecret: os.Getenv(CLIENTSECRET),
			Endpoint:     facebook.Endpoint,
			Scopes:       []string{"user_photos", "user_videos"},
			RedirectURL:  "http://localhost:8080/callback",
		}
	}
}

func main() {
	var input string
	fmt.Print("Enter the provider (google/facebook): ")
	fmt.Scanln(&input)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	oauthConfig = set_oauth(input)

	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", authURL)

	http.HandleFunc("/callback", callback)
	http.HandleFunc("/", redirect)

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	if code == "" {
		http.Error(w, "No code in the request", http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get(URL)
	if err != nil {
		http.Error(w, "Failed to get photos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Photos: %v\n", result)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
}
