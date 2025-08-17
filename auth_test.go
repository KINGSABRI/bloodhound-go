package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupTestServer creates a mock BloodHound API server for testing.
func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Login Endpoint
	mux.HandleFunc("/api/v2/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// A realistic JSON response for a successful login
		w.Write([]byte(`{
			"data": {
				"session_token": "test-session-token",
				"user_id": "12345",
				"user_dn": "CN=Test User,CN=Users,DC=CORP,DC=LOCAL"
			}
		}`))
	})

	// Mock Logout Endpoint
	mux.HandleFunc("/api/v2/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Check for auth header
		if r.Header.Get("Authorization") != "Bearer test-session-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return httptest.NewServer(mux)
}

func TestClient_Login(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	err = client.Login("testuser", "testpass")
	if err != nil {
		t.Errorf("Login failed unexpectedly: %v", err)
	}

	if client.GetToken() != "test-session-token" {
		t.Errorf("Expected token 'test-session-token', got '%s'", client.GetToken())
	}
}

func TestClient_Logout(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Manually set a token to test logout
	client.SetToken("test-session-token")

	err = client.Logout()
	if err != nil {
		t.Errorf("Logout failed unexpectedly: %v", err)
	}

	if client.GetToken() != "" {
		t.Errorf("Expected token to be empty after logout, got '%s'", client.GetToken())
	}
}