package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	baseURL     = "https://api.cloudflare.com/client/v4/accounts"
	defaultURLs = "google.com,facebook.com"
)

type Environment struct {
	APIKey    string
	Email     string
	AccountID string
	PoolID    string
}

type Origin struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type RequestBody struct {
	Origins []Origin `json:"origins"`
	Name    string   `json:"name"`
}

func initEnvs() (Environment, error) {
	env := Environment{
		APIKey:    os.Getenv("API_KEY"),
		Email:     os.Getenv("EMAIL"),
		AccountID: os.Getenv("ACCOUNT_ID"),
		PoolID:    os.Getenv("POOL_ID"),
	}

	if env.APIKey == "" || env.Email == "" || env.AccountID == "" || env.PoolID == "" {
		return Environment{}, fmt.Errorf("missing required environment variables")
	}

	return env, nil
}

func set(env Environment, urls string) error {
	log.Printf("Calling Cloudflare API with URLs: %s", urls)

	payload := CreatePayloadFromURLs(urls)
	jsonData, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	log.Printf("JSON Data: %s", string(jsonData))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut,
		fmt.Sprintf("%s/%s/load_balancers/pools/%s", baseURL, env.AccountID, env.PoolID),
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Email", env.Email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", env.APIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Response Status: %s", resp.Status)
	log.Printf("Response Headers: %v", resp.Header)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("Response Body: %s", string(body))
	return nil
}

func main() {
	env, err := initEnvs()
	if err != nil {
		log.Fatalf("Failed to initialize environment: %v", err)
	}

	urls := os.Getenv("URLS")
	if urls == "" {
		urls = defaultURLs
	}

	if err := set(env, urls); err != nil {
		log.Fatalf("Failed to set URLs: %v", err)
	}
}
