package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetADUser fetches a single AD user by their Object ID (SID).
func (c *Client) GetADUser(objectID string) (*ADUser, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID)
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
		return nil, fmt.Errorf("get ad user failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props                 ADUser `json:"props"`
			AdminRights           int    `json:"adminRights"`
			ConstrainedDelegation int    `json:"constrainedDelegation"`
			Controllables         int    `json:"controllables"`
			Controllers           int    `json:"controllers"`
			DCOMRights            int    `json:"dcomRights"`
			GroupMembership       int    `json:"groupMembership"`
			PSRemoteRights        int    `json:"psRemoteRights"`
			RDPRights             int    `json:"rdpRights"`
			Sessions              int    `json:"sessions"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get ad user response: %w", err)
	}

	user := response.Data.Props
	user.ObjectType = "User"
	user.AdminRights = response.Data.AdminRights
	user.ConstrainedDelegation = response.Data.ConstrainedDelegation
	user.Controllables = response.Data.Controllables
	user.Controllers = response.Data.Controllers
	user.DCOMRights = response.Data.DCOMRights
	user.GroupMembership = response.Data.GroupMembership
	user.PSRemoteRights = response.Data.PSRemoteRights
	user.RDPRights = response.Data.RDPRights
	user.Sessions = response.Data.Sessions
	return &user, nil
}

// GetADUserByName fetches a single AD user by their name.
func (c *Client) GetADUserByName(userName string) (*ADUser, error) {
	searchResponse, err := c.Search(userName, "User")
	if err != nil {
		return nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, userName) {
			return c.GetADUser(result.ObjectID)
		}
	}

	return nil, fmt.Errorf("user not found: %s", userName)
}

// GetADUserAdminRights fetches the admin rights for a given AD user.
func (c *Client) GetADUserAdminRights(objectID string) ([]EntityAdmin, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/admin-rights")
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

// GetADUserSessions fetches the sessions for a given AD user.
func (c *Client) GetADUserSessions(objectID string) ([]Session, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/sessions")
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

// GetADUserRDPRights fetches the RDP rights for a given AD user.
func (c *Client) GetADUserRDPRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/rdp-rights")
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

// GetADUserDCOMRights fetches the DCOM rights for a given AD user.
func (c *Client) GetADUserDCOMRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/dcom-rights")
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

// GetADUserPSRemoteRights fetches the PSRemote rights for a given AD user.
func (c *Client) GetADUserPSRemoteRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/ps-remote-rights")
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

// GetADUserSQLAdminRights fetches the SQL admin rights for a given AD user.
func (c *Client) GetADUserSQLAdminRights(objectID string) ([]Privilege, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/sql-admin-rights")
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

// GetADUserConstrainedDelegationRights fetches the constrained delegation rights for a given AD user.
func (c *Client) GetADUserConstrainedDelegationRights(objectID string) ([]ConstrainedDelegation, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/constrained-delegation-rights")
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

// GetADUserGroupMembership fetches the group membership for a given AD user.
func (c *Client) GetADUserGroupMembership(objectID string) ([]GroupMembership, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/memberships")
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

// GetADUserControllers fetches the controllers of a given AD user.
func (c *Client) GetADUserControllers(objectID string) ([]Controller, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/controllers")
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

// GetADUserControllables fetches the controllables of a given AD user.
func (c *Client) GetADUserControllables(objectID string) ([]Controllable, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/controllables")
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