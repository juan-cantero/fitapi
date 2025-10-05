package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	// Check for --json flag for machine-readable output
	jsonOutput := false
	if len(os.Args) > 1 && os.Args[1] == "--json" {
		jsonOutput = true
	}

	// Get email and password from args or use defaults
	email := "test@example.com"
	password := "test123456"

	if len(os.Args) > 2 && !jsonOutput {
		email = os.Args[1]
		password = os.Args[2]
	} else if len(os.Args) > 3 && jsonOutput {
		email = os.Args[2]
		password = os.Args[3]
	}

	// Try to sign in (if user exists)
	token, err := signIn(supabaseURL, supabaseKey, email, password)
	if err != nil {
		if !jsonOutput {
			fmt.Fprintln(os.Stderr, "Sign in failed, trying to sign up...")
		}
		// If sign in fails, try to sign up (create user)
		token, err = signUp(supabaseURL, supabaseKey, email, password)
		if err != nil {
			log.Fatalf("Sign up failed: %v", err)
		}
		if !jsonOutput {
			fmt.Fprintln(os.Stderr, "✅ User created successfully!")
		}
	}

	// Output format
	if jsonOutput {
		// Machine-readable JSON output
		output := map[string]interface{}{
			"access_token": token.AccessToken,
			"expires_in":   token.ExpiresIn,
			"expires_at":   calculateExpiresAt(token.ExpiresIn),
			"user_id":      token.User.ID,
			"email":        token.User.Email,
		}
		jsonData, _ := json.Marshal(output)
		fmt.Println(string(jsonData))
	} else {
		// Human-readable output
		fmt.Println("\n🎉 Authentication successful!")
		fmt.Println("\n📋 Copy this token for testing:")
		fmt.Println("─────────────────────────────────────────────────────────")
		fmt.Println(token.AccessToken)
		fmt.Println("─────────────────────────────────────────────────────────")
		fmt.Printf("\n👤 User ID: %s\n", token.User.ID)
		fmt.Printf("📧 Email: %s\n", token.User.Email)
		fmt.Printf("⏰ Expires in: %d seconds\n", token.ExpiresIn)
		fmt.Println("\n💡 Usage:")
		fmt.Println("curl http://localhost:8080/api/exercises \\")
		fmt.Printf("  -H 'Authorization: Bearer %s'\n", token.AccessToken)
	}
}

func calculateExpiresAt(expiresIn int) int64 {
	return time.Now().Unix() + int64(expiresIn)
}

func signIn(supabaseURL, apiKey, email, password string) (*SignInResponse, error) {
	url := fmt.Sprintf("%s/auth/v1/token?grant_type=password", supabaseURL)

	reqBody := SignInRequest{
		Email:    email,
		Password: password,
	}

	return makeAuthRequest(url, apiKey, reqBody)
}

func signUp(supabaseURL, apiKey, email, password string) (*SignInResponse, error) {
	url := fmt.Sprintf("%s/auth/v1/signup", supabaseURL)

	reqBody := SignInRequest{
		Email:    email,
		Password: password,
	}

	return makeAuthRequest(url, apiKey, reqBody)
}

func makeAuthRequest(url, apiKey string, reqBody SignInRequest) (*SignInResponse, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result SignInResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
