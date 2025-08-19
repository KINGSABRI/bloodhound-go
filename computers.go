package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetComputer fetches a single computer by its Object ID (SID).
func (c *Client) GetComputer(objectID string) (*Computer, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get computer failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props          Computer `json:"props"`
			AdminRights    int      `json:"adminRights"`
			Controllables  int      `json:"controllables"`
			Controllers    int      `json:"controllers"`
			DCOMRights     int      `json:"dcomRights"`
			PSRemoteRights int      `json:"psRemoteRights"`
			RDPRights      int      `json:"rdpRights"`
			Sessions       int      `json:"sessions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get computer response: %w", err)
	}

	computer := response.Data.Props
	computer.ObjectType = "Computer"
	computer.AdminRights = response.Data.AdminRights
	computer.Controllables = response.Data.Controllables
	computer.Controllers = response.Data.Controllers
	computer.DCOMRights = response.Data.DCOMRights
	computer.PSRemoteRights = response.Data.PSRemoteRights
	computer.RDPRights = response.Data.RDPRights
	computer.Sessions = response.Data.Sessions

	return &computer, nil
}

// GetComputerByName fetches a single computer by its name.
func (c *Client) GetComputerByName(computerName string) (*Computer, error) {
	searchResponse, err := c.Search(computerName, "Computer")
	if err != nil {
		return nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, computerName) {
			return c.GetComputer(result.ObjectID)
		}
	}

	return nil, fmt.Errorf("computer not found: %s", computerName)
}

// GetComputerAdmins fetches the list of principals with admin rights to a given computer.
func (c *Client) GetComputerAdmins(objectID string) ([]EntityAdmin, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/admin-rights")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse EntityAdminsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []EntityAdmin
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []EntityAdmin{}, nil
	}
	return finalResponse, nil
}

// GetComputerSessions fetches the user sessions on a given computer.
func (c *Client) GetComputerSessions(objectID string) ([]Session, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/sessions")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse SessionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Session
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Session{}, nil
	}
	return finalResponse, nil
}

// GetComputerRDPUsers fetches the principals with RDP rights to a given computer.
func (c *Client) GetComputerRDPUsers(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/rdp")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse PrivilegesResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Privilege
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Privilege{}, nil
	}
	return finalResponse, nil
}

// GetComputerDCOMUsers fetches the principals with DCOM rights to a given computer.
func (c *Client) GetComputerDCOMUsers(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/dcom")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse PrivilegesResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Privilege
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Privilege{}, nil
	}
	return finalResponse, nil
}

// GetComputerPSRemoteUsers fetches the principals with PSRemote rights to a given computer.
func (c *Client) GetComputerPSRemoteUsers(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/psremote")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse PrivilegesResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Privilege
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Privilege{}, nil
	}
	return finalResponse, nil
}

// GetComputerSQLAdmins fetches the principals with SQL admin rights to a given computer.
func (c *Client) GetComputerSQLAdmins(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/sql-admins")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse PrivilegesResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Privilege
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Privilege{}, nil
	}
	return finalResponse, nil
}

// GetComputerConstrainedDelegation fetches the constrained delegation privileges for a given computer.
func (c *Client) GetComputerConstrainedDelegation(objectID string) ([]ConstrainedDelegation, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/constrained-delegation-rights")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse ConstrainedDelegationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []ConstrainedDelegation
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []ConstrainedDelegation{}, nil
	}
	return finalResponse, nil
}

// GetComputerControllers fetches the controllers of a given computer.
func (c *Client) GetComputerControllers(objectID string) ([]Controller, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/controllers")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse ControllersResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Controller
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Controller{}, nil
	}
	return finalResponse, nil
}

// GetComputerControllables fetches the controllables of a given computer.
func (c *Client) GetComputerControllables(objectID string) ([]Controllable, error) {
	url := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/controllables")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse ControllablesResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Controllable
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Controllable{}, nil
	}
	return finalResponse, nil
}
