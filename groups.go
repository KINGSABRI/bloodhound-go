package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetGroup fetches a single group by its Object ID (SID).
func (c *Client) GetGroup(objectID string) (*Group, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID)
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
		return nil, fmt.Errorf("get group failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props          Group `json:"props"`
			AdminRights    int   `json:"adminRights"`
			Controllables  int   `json:"controllables"`
			Controllers    int   `json:"controllers"`
			DCOMRights     int   `json:"dcomRights"`
			Members        int   `json:"members"`
			PSRemoteRights int   `json:"psRemoteRights"`
			RDPRights      int   `json:"rdpRights"`
			Sessions       int   `json:"sessions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get group response: %w", err)
	}
	group := response.Data.Props
	group.ObjectType = "Group"
	group.AdminRights = response.Data.AdminRights
	group.Controllables = response.Data.Controllables
	group.Controllers = response.Data.Controllers
	group.DCOMRights = response.Data.DCOMRights
	group.Members = response.Data.Members
	group.PSRemoteRights = response.Data.PSRemoteRights
	group.RDPRights = response.Data.RDPRights
	group.Sessions = response.Data.Sessions
	return &group, nil
}

// GetGroupByName fetches a single group by its name.
func (c *Client) GetGroupByName(groupName string) (*Group, error) {
	searchResponse, err := c.Search(groupName, "Group", 0)
	if err != nil {
		return nil, err
	}

	if len(searchResponse.Data) == 0 && strings.Contains(groupName, "@") {
		samAccountName := strings.Split(groupName, "@")[0]
		searchResponse, err = c.Search(samAccountName, "Group", 0)
		if err != nil {
			return nil, err
		}
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, groupName) || (strings.Contains(groupName, "@") && strings.EqualFold(result.Name, strings.Split(groupName, "@")[0])) {
			return c.GetGroup(result.ObjectID)
		}
	}

	if len(searchResponse.Data) == 1 {
		return c.GetGroup(searchResponse.Data[0].ObjectID)
	}

	return nil, fmt.Errorf("group not found: %s", groupName)
}

// GetGroupMembers fetches the members of a given group.
func (c *Client) GetGroupMembers(objectID string, limit int) (GroupMembershipsResponse, error) {
	var rawResponse GroupMembershipsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/members")
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

// GetGroupMemberships fetches the group memberships for a given group.
func (c *Client) GetGroupMemberships(objectID string, limit int) (GroupMembershipsResponse, error) {
	var rawResponse GroupMembershipsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/memberships")
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

// GetGroupControllers fetches the controllers of a given group.
func (c *Client) GetGroupControllers(objectID string, limit int) (ControllersResponse, error) {
	var rawResponse ControllersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/controllers")
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

// GetGroupControllables fetches the controllables of a given group.
func (c *Client) GetGroupControllables(objectID string, limit int) (ControllablesResponse, error) {
	var rawResponse ControllablesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/controllables")
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

// GetGroupDCOMRights fetches principals with DCOM rights on the group.
func (c *Client) GetGroupDCOMRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/dcom-rights")
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

// GetGroupPSRemoteRights fetches principals with PSRemote rights on the group.
func (c *Client) GetGroupPSRemoteRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/ps-remote-rights")
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

// GetGroupRDPRights fetches principals with RDP rights on the group.
func (c *Client) GetGroupRDPRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/rdp-rights")
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

// GetGroupSessions fetches sessions on the group.
func (c *Client) GetGroupSessions(objectID string, limit int) (SessionsResponse, error) {
	var rawResponse SessionsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/sessions")
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
