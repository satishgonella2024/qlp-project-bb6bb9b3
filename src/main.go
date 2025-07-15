package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/login", loginHandler)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := server.ListenAndServeTLS(
		os.Getenv("TLS_CERT_FILE"),
		os.Getenv("TLS_KEY_FILE"),
	)
	if err != nil {
		panic(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate input
	if username == "" || password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// Fetch user from database and compare password
	hashedPassword, err := fetchHashedPassword(username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate and set secure session token
	sessionToken := generateSecureSessionToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

func fetchHashedPassword(username string) (string, error) {
	// Fetch hashed password from database
	// This is a placeholder, replace with actual database logic
	return "$2a$12$...", nil
}

func generateSecureSessionToken() string {
	// Generate a secure session token
	// This is a placeholder, replace with actual token generation logic
	return "secure_session_token"
}