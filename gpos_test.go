package bloodhound

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupGPOsTestServer creates a mock server with GPO-related endpoints.
func setupGPOsTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock Get GPO Endpoint
	mux.HandleFunc("/api/v2/gpos/D906C9B5-20EB-4111-98B2-6EF9CC7094B1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"base": {
					"objectid": "D906C9B5-20EB-4111-98B2-6EF9CC7094B1",
					"name": "DEFAULT DOMAIN POLICY@CORP.LOCAL",
					"distinguishedname": "CN={D906C9B5-20EB-4111-98B2-6EF9CC7094B1},CN=Policies,CN=System,DC=CORP,DC=LOCAL",
					"system_tags": "admin_tier_0",
					"whencreated": 1533502700,
					"type": "GPO"
				},
				"props": {}
			}
		}`))
	})

	// Mock Get GPO Controllers Endpoint
	mux.HandleFunc("/api/v2/gpos/D906C9B5-20EB-4111-98B2-6EF9CC7094B1/controllers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": [
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

func TestClient_GetGPO(t *testing.T) {
	server := setupGPOsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	_, base, err := client.GetGPO("D906C9B5-20EB-4111-98B2-6EF9CC7094B1")
	if err != nil {
		t.Fatalf("GetGPO failed unexpectedly: %v", err)
	}

	if base.Name != "DEFAULT DOMAIN POLICY@CORP.LOCAL" {
		t.Errorf("Expected GPO name 'DEFAULT DOMAIN POLICY@CORP.LOCAL', got '%s'", base.Name)
	}
}

func TestClient_GetGPOControllers(t *testing.T) {
	server := setupGPOsTestServer()
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	client.SetToken("test-token")

	controllers, err := client.GetGPOControllers("D906C9B5-20EB-4111-98B2-6EF9CC7094B1")
	if err != nil {
		t.Fatalf("GetGPOControllers failed unexpectedly: %v", err)
	}

	if len(controllers) != 1 {
		t.Errorf("Expected 1 controller, got %d", len(controllers))
	}
	if controllers[0].ObjectType != "Group" {
		t.Errorf("Expected controller to be a Group, got %s", controllers[0].ObjectType)
	}
}
