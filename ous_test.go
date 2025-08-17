package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupOUsTestServer creates a mock server with OU-related endpoints.
func setupOUsTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Get OU Endpoint
	mux.HandleFunc("/api/v2/ous/OU-123-456", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"base": {
					"objectid": "OU-123-456",
					"name": "DOMAIN CONTROLLERS@CORP.LOCAL",
					"distinguishedname": "OU=Domain Controllers,DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1533502700,
					"type": "OU"
				},
				"props": {}
			}
		}`))
	})

	// Mock Get OU Computers Endpoint
	mux.HandleFunc("/api/v2/ous/OU-123-456/computers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
				{
					"objectid": "S-1-5-21-123-456-1103",
					"name": "DC01.CORP.LOCAL",
					"distinguishedname": "CN=DC01,OU=Domain Controllers,DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1609459200,
					"type": "Computer",
					"operatingsystem": "Windows Server 2019",
					"enabled": true,
					"haslaps": true,
					"lastseen": 1672531200
				}
			]
		}`))
	})

	return httptest.NewServer(mux)
}

func TestClient_GetOU(t *testing.T) {
	server := setupOUsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	_, base, err := client.GetOU("OU-123-456")
	if err != nil {
		t.Fatalf("GetOU failed unexpectedly: %v", err)
	}

	if base.Name != "DOMAIN CONTROLLERS@CORP.LOCAL" {
		t.Errorf("Expected OU name 'DOMAIN CONTROLLERS@CORP.LOCAL', got '%s'", base.Name)
	}
}

func TestClient_GetOUComputers(t *testing.T) {
	server := setupOUsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	computers, err := client.GetOUComputers("OU-123-456")
	if err != nil {
		t.Fatalf("GetOUComputers failed unexpectedly: %v", err)
	}

	if len(computers) != 1 {
		t.Errorf("Expected 1 computer, got %d", len(computers))
	}
	if computers[0].Name != "DC01.CORP.LOCAL" {
		t.Errorf("Expected computer name 'DC01.CORP.LOCAL', got %s", computers[0].Name)
	}
}
