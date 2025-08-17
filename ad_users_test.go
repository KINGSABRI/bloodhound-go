package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupADUsersTestServer creates a mock server with AD user-related endpoints.
func setupADUsersTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Get User Endpoint
	mux.HandleFunc("/api/v2/users/S-1-5-21-123-456-500", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"base": {
					"objectid": "S-1-5-21-123-456-500",
					"name": "ADMINISTRATOR@CORP.LOCAL",
					"distinguishedname": "CN=Administrator,CN=Users,DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1533502700,
					"type": "User"
				},
				"props": {
					"enabled": true,
					"hassidhistory": false,
					"admincount": true
				}
			}
		}`))
	})

	// Mock Get User Memberships Endpoint
	mux.HandleFunc("/api/v2/users/S-1-5-21-123-456-500/memberships", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
				{
					"name": "DOMAIN ADMINS@CORP.LOCAL",
					"object_id": "S-1-5-21-123-456-512",
					"object_type": "Group",
					"is_tier_zero": true
				},
				{
					"name": "ENTERPRISE ADMINS@CORP.LOCAL",
					"object_id": "S-1-5-21-123-456-519",
					"object_type": "Group",
					"is_tier_zero": true
				}
			]
		}`))
	})

	return httptest.NewServer(mux)
}

func TestClient_GetADUser(t *testing.T) {
	server := setupADUsersTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	user, base, err := client.GetADUser("S-1-5-21-123-456-500")
	if err != nil {
		t.Fatalf("GetADUser failed unexpectedly: %v", err)
	}

	if base.Name != "ADMINISTRATOR@CORP.LOCAL" {
		t.Errorf("Expected user name 'ADMINISTRATOR@CORP.LOCAL', got '%s'", base.Name)
	}
	if !user.IsAdmin {
		t.Error("Expected IsAdmin to be true, but it was false")
	}
}

func TestClient_GetADUserGroupMembership(t *testing.T) {
	server := setupADUsersTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	memberships, err := client.GetADUserGroupMembership("S-1-5-21-123-456-500")
	if err != nil {
		t.Fatalf("GetADUserGroupMembership failed unexpectedly: %v", err)
	}

	if len(memberships) != 2 {
		t.Errorf("Expected 2 group memberships, got %d", len(memberships))
	}
	if memberships[1].Name != "ENTERPRISE ADMINS@CORP.LOCAL" {
		t.Errorf("Expected second group to be 'ENTERPRISE ADMINS@CORP.LOCAL', got %s", memberships[1].Name)
	}
}
