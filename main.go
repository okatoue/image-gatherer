package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConfig = &oauth2.Config {
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint:     google.Endpoint,
	Scopes:       []string{"https://www.googleapis.com/auth/photoslibrary.readonly"},
	RedirectURL:  "http://localhost:8080/callback",
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

	resp, err := client.Get("https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=10")
	
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

func main() {
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	
	fmt.Printf("Visit the URL for the auth dialog: %v\n", authURL)

	http.HandleFunc("/callback", callback)

	http.HandleFunc("/", redirect)
	
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}