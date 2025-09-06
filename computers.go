package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetComputer fetches a single computer by its Object ID (SID).
func (c *Client) GetComputer(objectID string) (*Computer, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
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
			Props                   Computer `json:"props"`
			AdminRights             int      `json:"adminRights"`
			AdminUsers              int      `json:"adminUsers"`
			ConstrainedPrivs        int      `json:"constrainedPrivs"`
			ConstrainedUsers        int      `json:"constrainedUsers"`
			Controllables           int      `json:"controllables"`
			Controllers             int      `json:"controllers"`
			DCOMRights              int      `json:"dcomRights"`
			DCOMUsers               int      `json:"dcomUsers"`
			GPOs                    int      `json:"gpos"`
			GroupMembership         int      `json:"groupMembership"`
			PSRemoteRights          int      `json:"psRemoteRights"`
			PSRemoteUsers           int      `json:"psRemoteUsers"`
			RDPRights               int      `json:"rdpRights"`
			Sessions                int      `json:"sessions"`
			SQLAdminUsers           int      `json:"sqlAdminUsers"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get computer response: %w", err)
	}

	computer := response.Data.Props
	computer.ObjectType = "Computer"
	computer.AdminRights = response.Data.AdminRights
	computer.AdminUsers = response.Data.AdminUsers
	computer.ConstrainedPrivs = response.Data.ConstrainedPrivs
	computer.ConstrainedUsers = response.Data.ConstrainedUsers
	computer.Controllables = response.Data.Controllables
	computer.Controllers = response.Data.Controllers
	computer.DCOMRights = response.Data.DCOMRights
	computer.DCOMUsers = response.Data.DCOMUsers
	computer.GPOs = response.Data.GPOs
	computer.GroupMembership = response.Data.GroupMembership
	computer.PSRemoteRights = response.Data.PSRemoteRights
	computer.PSRemoteUsers = response.Data.PSRemoteUsers
	computer.RDPRights = response.Data.RDPRights
	computer.Sessions = response.Data.Sessions
	computer.SQLAdminUsers = response.Data.SQLAdminUsers

	return &computer, nil
}

// GetComputerByName fetches a single computer by its name.
func (c *Client) GetComputerByName(computerName string) (*Computer, error) {
	searchResponse, err := c.Search(computerName, "Computer", 0)
	if err != nil {
		return nil, err
	}

	if len(searchResponse.Data) == 0 && strings.Contains(computerName, "@") {
		samAccountName := strings.Split(computerName, "@")[0]
		searchResponse, err = c.Search(samAccountName, "Computer", 0)
		if err != nil {
			return nil, err
		}
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, computerName) || (strings.Contains(computerName, "@") && strings.EqualFold(result.Name, strings.Split(computerName, "@")[0])) {
			return c.GetComputer(result.ObjectID)
		}
	}

	if len(searchResponse.Data) == 1 {
		return c.GetComputer(searchResponse.Data[0].ObjectID)
	}

	return nil, fmt.Errorf("computer not found: %s", computerName)
}

// GetComputerAdmins fetches the list of principals with admin rights to a given computer.
func (c *Client) GetComputerAdmins(objectID string, limit int) (EntityAdminsResponse, error) {
	var rawResponse EntityAdminsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/admin-rights")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerSessions fetches the user sessions on a given computer.
func (c *Client) GetComputerSessions(objectID string, limit int) (SessionsResponse, error) {
	var rawResponse SessionsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/sessions")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerRDPUsers fetches the principals with RDP rights to a given computer.
func (c *Client) GetComputerRDPUsers(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/rdp-rights")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerDCOMUsers fetches the principals with DCOM rights to a given computer.
func (c *Client) GetComputerDCOMUsers(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/dcom-rights")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerPSRemoteUsers fetches the principals with PSRemote rights to a given computer.
func (c *Client) GetComputerPSRemoteUsers(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/ps-remote-rights")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerSQLAdmins fetches the principals with SQL admin rights to a given computer.
func (c *Client) GetComputerSQLAdmins(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/sql-admins")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerConstrainedDelegation fetches the constrained delegation privileges for a given computer.
func (c *Client) GetComputerConstrainedDelegation(objectID string, limit int) (ConstrainedDelegationsResponse, error) {
	var rawResponse ConstrainedDelegationsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/constrained-delegation-rights")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerControllers fetches the controllers of a given computer.
func (c *Client) GetComputerControllers(objectID string, limit int) (ControllersResponse, error) {
	var rawResponse ControllersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/controllers")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerMemberships fetches the group memberships for a given computer.
func (c *Client) GetComputerMemberships(objectID string, limit int) (GroupMembershipsResponse, error) {
	var rawResponse GroupMembershipsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/group-membership")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetComputerControllables fetches the controllables of a given computer.
func (c *Client) GetComputerControllables(objectID string, limit int) (ControllablesResponse, error) {
	var rawResponse ControllablesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/computers/", objectID, "/controllables")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}