package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
			GPOs                  int    `json:"gpos"`
			PSRemoteRights        int    `json:"psRemoteRights"`
			RDPRights             int    `json:"rdpRights"`
			Sessions              int    `json:"sessions"`
			SQLAdminRights        int    `json:"sqlAdminRights"`
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
	user.GPOs = response.Data.GPOs
	user.PSRemoteRights = response.Data.PSRemoteRights
	user.RDPRights = response.Data.RDPRights
	user.Sessions = response.Data.Sessions
	user.SQLAdminRights = response.Data.SQLAdminRights
	return &user, nil
}

// GetADUserByName fetches a single AD user by their name.
func (c *Client) GetADUserByName(userName string) (*ADUser, error) {
	// First, try searching with the provided username directly
	searchResponse, err := c.Search(userName, "User", 0)
	if err != nil {
		return nil, err
	}

	// If the initial search yields no results and the username contains an "@",
	// try searching for the part before the "@" as a fallback.
	if len(searchResponse.Data) == 0 && strings.Contains(userName, "@") {
		samAccountName := strings.Split(userName, "@")[0]
		searchResponse, err = c.Search(samAccountName, "User", 0)
		if err != nil {
			return nil, err
		}
	}

	// Now, process the results from the successful search
	for _, result := range searchResponse.Data {
		// Use EqualFold for case-insensitive comparison on the result's name
		if strings.EqualFold(result.Name, userName) || (strings.Contains(userName, "@") && strings.EqualFold(result.Name, strings.Split(userName, "@")[0])) {
			return c.GetADUser(result.ObjectID)
		}
	}

	// If we still haven't found a match, try a direct hit on the first result if there's only one.
	// This handles cases where the name in BH is slightly different (e.g. UPN vs SAM)
	if len(searchResponse.Data) == 1 {
		return c.GetADUser(searchResponse.Data[0].ObjectID)
	}


	return nil, fmt.Errorf("user not found: %s", userName)
}

// GetADUserAdminRights fetches the admin rights for a given AD user.
func (c *Client) GetADUserAdminRights(objectID string, limit int) (EntityAdminsResponse, error) {
	var rawResponse EntityAdminsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/admin-rights")
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

// GetADUserSessions fetches the sessions for a given AD user.
func (c *Client) GetADUserSessions(objectID string, limit int) (SessionsResponse, error) {
	var rawResponse SessionsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/sessions")
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

// GetADUserRDPRights fetches the RDP rights for a given AD user.
func (c *Client) GetADUserRDPRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/rdp-rights")
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

// GetADUserDCOMRights fetches the DCOM rights for a given AD user.
func (c *Client) GetADUserDCOMRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/dcom-rights")
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

// GetADUserPSRemoteRights fetches the PSRemote rights for a given AD user.
func (c *Client) GetADUserPSRemoteRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/ps-remote-rights")
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

// GetADUserSQLAdminRights fetches the SQL admin rights for a given AD user.
func (c *Client) GetADUserSQLAdminRights(objectID string, limit int) (PrivilegesResponse, error) {
	var rawResponse PrivilegesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/sql-admin-rights")
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

// GetADUserConstrainedDelegationRights fetches the constrained delegation rights for a given AD user.
func (c *Client) GetADUserConstrainedDelegationRights(objectID string, limit int) (ConstrainedDelegationsResponse, error) {
	var rawResponse ConstrainedDelegationsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/constrained-delegation-rights")
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

// GetADUserGroupMembership fetches the group membership for a given AD user.
func (c *Client) GetADUserGroupMembership(objectID string, limit int) (GroupMembershipsResponse, error) {
	var rawResponse GroupMembershipsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/memberships")
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

// GetADUserControllers fetches the controllers of a given AD user.
func (c *Client) GetADUserControllers(objectID string, limit int) (ControllersResponse, error) {
	var rawResponse ControllersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/controllers")
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

// GetADUserControllables fetches the controllables of a given AD user.
func (c *Client) GetADUserControllables(objectID string, limit int) (ControllablesResponse, error) {
	var rawResponse ControllablesResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/users/", objectID, "/controllables")
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

// ResolveUserIdentity takes a user identity (name or SID) and returns the SID.
// If a name is provided, it will be resolved via the search API.
func (c *Client) ResolveUserIdentity(identity string) (string, error) {
	// If the identity looks like a SID, return it directly.
	if strings.HasPrefix(strings.ToUpper(identity), "S-1-5-") {
		return identity, nil
	}

	// If not a SID, assume it's a name and search for it.
	searchResponse, err := c.Search(identity, "User", 0)
	if err != nil {
		return "", fmt.Errorf("search failed for user '%s': %w", identity, err)
	}

	if len(searchResponse.Data) == 0 {
		return "", fmt.Errorf("no user found matching name '%s'", identity)
	}

	// Prefer an exact (case-insensitive) match if there are multiple results.
	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, identity) {
			return result.ObjectID, nil
		}
	}

	// If no exact match was found and there are multiple results, it's ambiguous.
	if len(searchResponse.Data) > 1 {
		return "", fmt.Errorf("multiple users found matching name '%s'. Please be more specific or use the Object ID (SID)", identity)
	}

	// Only one result, so we'll use it.
	return searchResponse.Data[0].ObjectID, nil
}