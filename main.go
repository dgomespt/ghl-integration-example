package main

import (
	"context"
	"crypto/rand"
	"dgomespt/oauth2-example-app/pkg/response"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type GoHighLevel struct {
	Endpoint oauth2.Endpoint
}

var goHighLevel = GoHighLevel{
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://services.leadconnectorhq.com/oauth/token",
		TokenURL: "https://services.leadconnectorhq.com/oauth/token",
	},
}

func generateStateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("Failed to generate state token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

var oauthConfig *oauth2.Config
var stateToken string
var client *http.Client
var token *oauth2.Token

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       strings.Split(os.Getenv("OAUTH2_SCOPES"), ","),
		Endpoint:     goHighLevel.Endpoint,
	}

}

func generateAuthorizePageURL() string {

	params := []any{
		os.Getenv("OAUTH2_BASE_URL"),
		oauthConfig.RedirectURL,
		oauthConfig.ClientID,
		strings.Join(oauthConfig.Scopes, " "),
		stateToken,
	}

	return fmt.Sprintf("%s/oauth/chooselocation?response_type=code&redirect_uri=%s&client_id=%s&scope=%s&state=%s", params...)
}

func main() {

	loadConfig()

	stateToken = generateStateToken()

	t, err := loadToken()

	if err == nil {
		token = t
		log.Println("Using saved token")
		client = oauthConfig.Client(context.Background(), t)
	} else {
		log.Println("No saved token, login required")
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/me", handleGetMe)
	http.HandleFunc("/contacts", handleGetContacts)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><a href='/login'>Login</a></body></html>"))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	//url := oauthConfig.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)
	http.Redirect(w, r, generateAuthorizePageURL(), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate state token
	state := r.URL.Query().Get("state")
	if state != stateToken {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for an access token
	code := r.URL.Query().Get("code")
	tok, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		log.Printf("Error exchanging token: %v", err)
		return
	}

	// Save token to file
	if err := saveToken(tok); err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		log.Printf("Error saving token: %v", err)
		return
	}

	t, err := loadToken()
	if err != nil {
		token = t
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><p>Access token saved</p></html>"))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		log.Printf("Error reading webhook body: %v", err)
		return
	}
	r.Body.Close()

	log.Printf("Webhook received: %s", string(body))

	// Respond to the webhook request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}

func saveToken(token *oauth2.Token) error {
	file, err := os.Create("token.json")
	if err != nil {
		return fmt.Errorf("failed to create token file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		return fmt.Errorf("failed to encode token to file: %w", err)
	}
	return nil
}

func handleGetMe(w http.ResponseWriter, r *http.Request) {

	resp, err := client.Get("https://api.gohighlevel.com/v1/me")
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		log.Printf("Error fetching user info: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		log.Printf("Non-200 response fetching user info: %v", resp.Status)
		return
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		log.Printf("Error decoding user info: %v", err)
		return
	}

	// Display user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func loadToken() (*oauth2.Token, error) {
	file, err := os.Open("token.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open token file: %w", err)
	}
	defer file.Close()

	var token oauth2.Token
	if err := json.NewDecoder(file).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token from file: %w", err)
	}

	if token.Expiry.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return &token, nil
}

func handleGetContacts(w http.ResponseWriter, r *http.Request) {

	locationId := r.URL.Query().Get("locationId")

	url := "https://services.leadconnectorhq.com/contacts/search"

	data := map[string]any{
		"page":       1,
		"pageLimit":  10,
		"locationId": locationId,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	payload := strings.NewReader(string(jsonData))

	req, _ := http.NewRequest("POST", url, payload)

	if token == nil {
		http.Error(w, "No valid token found", http.StatusUnauthorized)
		log.Println("No valid token found")
		return
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Version", "2021-07-28")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		log.Printf("Error fetching contacts: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(resp)
	fmt.Println(string(body))

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		log.Printf("Non-200 response fetching contacts: %v", resp.Status)
		return
	}

	var contactsResponse response.ContactsResponse
	err = json.Unmarshal(body, &contactsResponse)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(contactsResponse)
	if err != nil {
		fmt.Println("Error displaying response:", err)
		return
	}
}
