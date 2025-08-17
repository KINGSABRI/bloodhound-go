package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupGroupsTestServer creates a mock server with group-related endpoints.
func setupGroupsTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Get Group Endpoint
	mux.HandleFunc("/api/v2/groups/S-1-5-21-123-456-512", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"base": {
					"objectid": "S-1-5-21-123-456-512",
					"name": "DOMAIN ADMINS@CORP.LOCAL",
					"distinguishedname": "CN=Domain Admins,CN=Users,DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1533502700,
					"type": "Group"
				},
				"props": {
					"admincount": true
				}
			}
		}`))
	})

	// Mock Get Group Members Endpoint
	mux.HandleFunc("/api/v2/groups/S-1-5-21-123-456-512/members", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
				{
					"name": "ADMINISTRATOR@CORP.LOCAL",
					"object_id": "S-1-5-21-123-456-500",
					"object_type": "User",
					"is_tier_zero": true
				}
			]
		}`))
	})

	return httptest.NewServer(mux)
}

func TestClient_GetGroup(t *testing.T) {
	server := setupGroupsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	group, base, err := client.GetGroup("S-1-5-21-123-456-512")
	if err != nil {
		t.Fatalf("GetGroup failed unexpectedly: %v", err)
	}

	if base.Name != "DOMAIN ADMINS@CORP.LOCAL" {
		t.Errorf("Expected group name 'DOMAIN ADMINS@CORP.LOCAL', got '%s'", base.Name)
	}
	if !group.IsAdmin {
		t.Error("Expected IsAdmin to be true, but it was false")
	}
}

func TestClient_GetGroupMembers(t *testing.T) {
	server := setupGroupsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	members, err := client.GetGroupMembers("S-1-5-21-123-456-512")
	if err != nil {
		t.Fatalf("GetGroupMembers failed unexpectedly: %v", err)
	}

	if len(members) != 1 {
		t.Errorf("Expected 1 group member, got %d", len(members))
	}
	if members[0].ObjectType != "User" {
		t.Errorf("Expected member to be a User, got %s", members[0].ObjectType)
	}
}
