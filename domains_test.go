package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupDomainsTestServer creates a mock server with domain-related endpoints.
func setupDomainsTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Get Domain Endpoint
	mux.HandleFunc("/api/v2/domains/S-1-5-21-123-456", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"base": {
					"objectid": "S-1-5-21-123-456",
					"name": "CORP.LOCAL",
					"distinguishedname": "DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1533502700,
					"type": "Domain"
				},
				"props": {}
			}
		}`))
	})

	// Mock Get Domain Controllers Endpoint
	mux.HandleFunc("/api/v2/domains/S-1-5-21-123-456/controllers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
				{
					"name": "DOMAIN ADMINS@CORP.LOCAL",
					"object_id": "S-1-5-21-123-456-512",
					"object_type": "Group",
					"is_tier_zero": true
				}
			]
		}`))
	})

	return httptest.NewServer(mux)
}

func TestClient_GetDomain(t *testing.T) {
	server := setupDomainsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	_, base, err := client.GetDomain("S-1-5-21-123-456")
	if err != nil {
		t.Fatalf("GetDomain failed unexpectedly: %v", err)
	}

	if base.Name != "CORP.LOCAL" {
		t.Errorf("Expected domain name 'CORP.LOCAL', got '%s'", base.Name)
	}
}

func TestClient_GetDomainControllers(t *testing.T) {
	server := setupDomainsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	controllers, err := client.GetDomainControllers("S-1-5-21-123-456")
	if err != nil {
		t.Fatalf("GetDomainControllers failed unexpectedly: %v", err)
	}

	if len(controllers) != 1 {
		t.Errorf("Expected 1 controller, got %d", len(controllers))
	}
	if controllers[0].ObjectType != "Group" {
		t.Errorf("Expected controller to be a Group, got %s", controllers[0].ObjectType)
	}
}
