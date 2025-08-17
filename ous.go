package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetOU fetches a single OU by its Object ID (SID).
func (c *Client) GetOU(objectID string) (*OU, *BaseEntity, error) {
	url := c.baseURL.JoinPath("/api/v2/ous/", objectID)
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
		return nil, nil, fmt.Errorf("get ou failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props      OU         `json:"props"`
			BaseEntity BaseEntity `json:"base"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode get ou response: %w", err)
	}
	response.Data.Props.ObjectType = "OU"
	return &response.Data.Props, &response.Data.BaseEntity, nil
}

// GetOUByName fetches a single OU by its name.
func (c *Client) GetOUByName(ouName string) (*OU, *BaseEntity, error) {
	searchResponse, err := c.Search(ouName, "OU")
	if err != nil {
		return nil, nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, ouName) {
			fullOU, baseEntity, err := c.GetOU(result.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			return fullOU, baseEntity, nil
		}
	}

	return nil, nil, fmt.Errorf("ou not found: %s", ouName)
}

// GetOUComputers fetches the computers in a given OU.
func (c *Client) GetOUComputers(objectID string) ([]Computer, error) {
	url := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/computers")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Computer
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Computer{}, nil
	}
	return finalResponse, nil
}

// GetOUUsers fetches the users in a given OU.
func (c *Client) GetOUUsers(objectID string) ([]ADUser, error) {
	url := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/users")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []ADUser
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []ADUser{}, nil
	}
	return finalResponse, nil
}

// GetOUGroups fetches the groups in a given OU.
func (c *Client) GetOUGroups(objectID string) ([]Group, error) {
	url := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/groups")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Group
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Group{}, nil
	}
	return finalResponse, nil
}

// GetOUControllers fetches the controllers of a given OU.
func (c *Client) GetOUControllers(objectID string) ([]Controller, error) {
	url := c.baseURL.JoinPath("/api/v2/ous/", objectID, "/controllers")
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
