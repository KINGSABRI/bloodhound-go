package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupSearchTestServer creates a mock server with a search endpoint.
func setupSearchTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v2/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Check for auth header
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		query := r.URL.Query().Get("q")
		searchType := r.URL.Query().Get("type")

		// Basic validation for the test
		if query == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "search query is required"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Return a different response based on the type parameter for testing
		if searchType == "User" {
			w.Write([]byte(`{
				"data": [
					{
						"objectid": "S-1-5-21-123-456-500",
						"type": "User",
						"name": "ADMINISTRATOR@CORP.LOCAL",
						"system_tags": "admin_tier_0"
					}
				]
			}`))
		} else {
			w.Write([]byte(`{
				"data": [
					{
						"objectid": "S-1-5-21-123-456-500",
						"type": "User",
						"name": "ADMINISTRATOR@CORP.LOCAL",
						"system_tags": "admin_tier_0"
					},
					{
						"objectid": "S-1-5-21-123-456-1103",
						"type": "Computer",
						"name": "DC01.CORP.LOCAL",
						"system_tags": "admin_tier_0"
					}
				]
			}`))
		}
	})

	return httptest.NewServer(mux)
}

func TestClient_Search(t *testing.T) {
	server := setupSearchTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token") // Set a dummy token for authenticated requests

	// Test case 1: General search without type
	resp, err := client.Search("admin", "")
	if err != nil {
		t.Fatalf("Search failed unexpectedly: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 search results, got %d", len(resp.Data))
	}
	if resp.Data[1].Type != "Computer" {
		t.Errorf("Expected second result to be a Computer, got %s", resp.Data[1].Type)
	}

	// Test case 2: Search with a specific type
	resp, err = client.Search("administrator", "User")
	if err != nil {
		t.Fatalf("Search with type failed unexpectedly: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 search result when filtering by type, got %d", len(resp.Data))
	}
	if resp.Data[0].Type != "User" {
		t.Errorf("Expected result to be a User, got %s", resp.Data[0].Type)
	}
}
