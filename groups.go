package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetGroup fetches a single group by its Object ID (SID).
func (c *Client) GetGroup(objectID string) (*Group, *BaseEntity, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get group failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props      Group      `json:"props"`
			BaseEntity BaseEntity `json:"base"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode get group response: %w", err)
	}
	response.Data.Props.ObjectType = "Group"
	return &response.Data.Props, &response.Data.BaseEntity, nil
}

// GetGroupByName fetches a single group by its name.
func (c *Client) GetGroupByName(groupName string) (*Group, *BaseEntity, error) {
	searchResponse, err := c.Search(groupName, "Group")
	if err != nil {
		return nil, nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, groupName) {
			fullGroup, baseEntity, err := c.GetGroup(result.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			return fullGroup, baseEntity, nil
		}
	}

	return nil, nil, fmt.Errorf("group not found: %s", groupName)
}

// GetGroupMembers fetches the members of a given group.
func (c *Client) GetGroupMembers(objectID string) ([]GroupMembership, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/members")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse GroupMembershipsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []GroupMembership
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []GroupMembership{}, nil
	}
	return finalResponse, nil
}

// GetGroupControllers fetches the controllers of a given group.
func (c *Client) GetGroupControllers(objectID string) ([]Controller, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/controllers")
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

// GetGroupControllables fetches the controllables of a given group.
func (c *Client) GetGroupControllables(objectID string) ([]Controllable, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/controllables")
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

// GetGroupDCOMRights fetches principals with DCOM rights on the group.
func (c *Client) GetGroupDCOMRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/dcom-rights")
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

// GetGroupPSRemoteRights fetches principals with PSRemote rights on the group.
func (c *Client) GetGroupPSRemoteRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/psremote-rights")
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

// GetGroupRDPRights fetches principals with RDP rights on the group.
func (c *Client) GetGroupRDPRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/rdp-rights")
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

// GetGroupSessions fetches sessions on the group.
func (c *Client) GetGroupSessions(objectID string) ([]Session, error) {
	url := c.baseURL.JoinPath("/api/v2/groups/", objectID, "/sessions")
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