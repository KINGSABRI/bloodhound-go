package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetOU fetches a single OU by its Object ID (GUID).
func (c *Client) GetOU(objectID string) (*OU, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/ous/", objectID)
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
		return nil, fmt.Errorf("get ou failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props     OU  `json:"props"`
			Computers int `json:"computers"`
			GPOs      int `json:"gpos"`
			Groups    int `json:"groups"`
			Users     int `json:"users"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get ou response: %w", err)
	}
	ou := response.Data.Props
	ou.ObjectType = "OU"
	ou.Computers = response.Data.Computers
	ou.GPOs = response.Data.GPOs
	ou.Groups = response.Data.Groups
	ou.Users = response.Data.Users
	return &ou, nil
}

// GetOUByName fetches a single OU by its name.
func (c *Client) GetOUByName(ouName string) (*OU, error) {
	searchResponse, err := c.Search(ouName, "OU", 0)
	if err != nil {
		return nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, ouName) {
			return c.GetOU(result.ObjectID)
		}
	}

	return nil, fmt.Errorf("ou not found: %s", ouName)
}

// GetOUGroups fetches the groups in a given OU.
func (c *Client) GetOUGroups(objectID string, limit int) (GroupsResponse, error) {
	var rawResponse GroupsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/groups")
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

// GetOUComputers fetches the computers in a given OU.
func (c *Client) GetOUComputers(objectID string, limit int) (ComputersResponse, error) {
	var rawResponse ComputersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/computers")
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

// GetOUUsers fetches the users in a given OU.
func (c *Client) GetOUUsers(objectID string, limit int) (UsersResponse, error) {
	var rawResponse UsersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/users")
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

// GetOuGPOs fetches the GPOs linked to a given OU.
func (c *Client) GetOuGPOs(objectID string, limit int) (GPOsResponse, error) {
	var rawResponse GPOsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/gpos")
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
