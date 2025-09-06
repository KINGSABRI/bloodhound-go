package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetContainer fetches a single container by its Object ID (GUID).
func (c *Client) GetContainer(objectID string) (*Container, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/containers/", objectID)
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
		return nil, fmt.Errorf("get container failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props Container `json:"props"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get container response: %w", err)
	}
	container := response.Data.Props
	container.ObjectType = "Container"
	return &container, nil
}

// GetContainerByName fetches a single container by its name.
func (c *Client) GetContainerByName(containerName string) (*Container, error) {
	searchResponse, err := c.Search(containerName, "Container", 0)
	if err != nil {
		return nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, containerName) {
			return c.GetContainer(result.ObjectID)
		}
	}

	return nil, fmt.Errorf("container not found: %s", containerName)
}

// GetContainerUsers fetches the users in a given container.
func (c *Client) GetContainerUsers(objectID string, limit int) (UsersResponse, error) {
	var rawResponse UsersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/containers/", objectID, "/users")
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

// GetContainerComputers fetches the computers in a given container.
func (c *Client) GetContainerComputers(objectID string, limit int) (ComputersResponse, error) {
	var rawResponse ComputersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/containers/", objectID, "/computers")
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

// GetContainerGroups fetches the groups in a given container.
func (c *Client) GetContainerGroups(objectID string, limit int) (GroupsResponse, error) {
	var rawResponse GroupsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/containers/", objectID, "/groups")
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

// GetContainerControllers fetches the controllers of a given container.
func (c *Client) GetContainerControllers(objectID string, limit int) (ControllersResponse, error) {
	var rawResponse ControllersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/containers/", objectID, "/controllers")
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
