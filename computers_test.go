package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupComputersTestServer creates a mock server with computer-related endpoints.
func setupComputersTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Get Computer Endpoint
	mux.HandleFunc("/api/v2/computers/S-1-5-21-123-456-1103", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"base": {
					"objectid": "S-1-5-21-123-456-1103",
					"name": "DC01.CORP.LOCAL",
					"distinguishedname": "CN=DC01,OU=Domain Controllers,DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1609459200,
					"type": "Computer"
				},
				"props": {
					"operatingsystem": "Windows Server 2019",
					"enabled": true,
					"haslaps": true,
					"lastseen": 1672531200
				}
			}
		}`))
	})

	// Mock Get Computer Admins Endpoint
	mux.HandleFunc("/api/v2/computers/S-1-5-21-123-456-1103/admins", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
				{
					"name": "ADMINISTRATOR@CORP.LOCAL",
					"object_id": "S-1-5-21-123-456-500",
					"object_type": "User",
					"is_tier_zero": true
				},
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

func TestClient_GetComputer(t *testing.T) {
	server := setupComputersTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	computer, base, err := client.GetComputer("S-1-5-21-123-456-1103")
	if err != nil {
		t.Fatalf("GetComputer failed unexpectedly: %v", err)
	}

	if base.Name != "DC01.CORP.LOCAL" {
		t.Errorf("Expected computer name 'DC01.CORP.LOCAL', got '%s'", base.Name)
	}
	if computer.OperatingSystem != "Windows Server 2019" {
		t.Errorf("Expected OS 'Windows Server 2019', got '%s'", computer.OperatingSystem)
	}
	if !computer.HasLAPS {
		t.Error("Expected HasLAPS to be true, but it was false")
	}
}

func TestClient_GetComputerAdmins(t *testing.T) {
	server := setupComputersTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	admins, err := client.GetComputerAdmins("S-1-5-21-123-456-1103")
	if err != nil {
		t.Fatalf("GetComputerAdmins failed unexpectedly: %v", err)
	}

	if len(admins) != 2 {
		t.Errorf("Expected 2 admins, got %d", len(admins))
	}
	if admins[1].ObjectType != "Group" {
		t.Errorf("Expected second admin to be a Group, got %s", admins[1].ObjectType)
	}
}
