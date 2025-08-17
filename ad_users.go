package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetADUser fetches a single AD user by their Object ID (SID).
func (c *Client) GetADUser(objectID string) (*ADUser, *BaseEntity, error) {
	url := c.baseURL.JoinPath("/api/v2/users/", objectID)
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
		return nil, nil, fmt.Errorf("get ad user failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props      ADUser     `json:"props"`
			BaseEntity BaseEntity `json:"base"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode get ad user response: %w", err)
	}
	response.Data.Props.ObjectType = "User"
	return &response.Data.Props, &response.Data.BaseEntity, nil
}

// GetADUserByName fetches a single AD user by their name.
func (c *Client) GetADUserByName(userName string) (*ADUser, *BaseEntity, error) {
	searchResponse, err := c.Search(userName, "User")
	if err != nil {
		return nil, nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, userName) {
			fullUser, baseEntity, err := c.GetADUser(result.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			return fullUser, baseEntity, nil
		}
	}

	return nil, nil, fmt.Errorf("user not found: %s", userName)
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
		// If it's not a list, it's probably a number (0), so we return an empty list
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
	url := c.baseURL.JoinPath("/api/v2/users/", objectID, "/psremote-rights")
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