package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Login authenticates to the BloodHound API using a username and password.
func (c *Client) Login(username, password string) error {
	loginRequest := LoginRequest{
		LoginMethod: "secret",
		Username:    username,
		Secret:      password,
	}
	payload, err := json.Marshal(loginRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	// For debugging, create a sanitized version of the request body
	sanitizedRequest := loginRequest
	if sanitizedRequest.Secret != "" {
		sanitizedRequest.Secret = "[REDACTED]"
	}
	sanitizedPayload, _ := json.Marshal(sanitizedRequest)

	// Build the request
	loginURL := c.baseURL.JoinPath("/api/v2/login")
	req, err := http.NewRequest(http.MethodPost, loginURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}
	// Set headers to match the successful curl request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip")

	// Execute the request using the special do method with sanitized body
	resp, err := c.do(req, sanitizedPayload)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status code: %d", resp.StatusCode)
	}

	var sessionResponse SessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sessionResponse); err != nil {
		return fmt.Errorf("failed to decode session response: %w", err)
	}

	c.setAuthToken(sessionResponse.Data.SessionToken)
	return nil
}

// Logout invalidates the current session token.
func (c *Client) Logout() error {
	logoutURL := c.baseURL.JoinPath("/api/v2/logout")
	req, err := c.newAuthenticatedRequest(http.MethodPost, logoutURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute logout request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status code: %d", resp.StatusCode)
	}

	c.setAuthToken("")
	return nil
}
