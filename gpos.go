package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetGPO fetches a single GPO by its Object ID (SID).
func (c *Client) GetGPO(objectID string) (*GPO, error) {
	url := c.baseURL.JoinPath("/api/v2/gpos/", objectID)
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
		return nil, fmt.Errorf("get gpo failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props GPO `json:"props"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get gpo response: %w", err)
	}
	gpo := response.Data.Props
	gpo.ObjectType = "GPO"
	return &gpo, nil
}

// GetGPOByName fetches a single GPO by its name.
func (c *Client) GetGPOByName(gpoName string) (*GPO, error) {
	searchResponse, err := c.Search(gpoName, "GPO")
	if err != nil {
		return nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, gpoName) {
			return c.GetGPO(result.ObjectID)
		}
	}

	return nil, fmt.Errorf("gpo not found: %s", gpoName)
}

// GetGPOControllers fetches the controllers of a given GPO.
func (c *Client) GetGPOControllers(objectID string) ([]EntityAdmin, error) {
	url := c.baseURL.JoinPath("/api/v2/gpos/", objectID, "/controllers")
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

// GetGPOComputers fetches the computers affected by a given GPO.
func (c *Client) GetGPOComputers(objectID string) ([]Computer, error) {
	url := c.baseURL.JoinPath("/api/v2/gpos/", objectID, "/computers")
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

// GetGPOUsers fetches the users affected by a given GPO.
func (c *Client) GetGPOUsers(objectID string) ([]ADUser, error) {
	url := c.baseURL.JoinPath("/api/v2/gpos/", objectID, "/users")
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

// GetGPOOUs fetches the OUs affected by a given GPO.
func (c *Client) GetGPOOUs(objectID string) ([]OU, error) {
	url := c.baseURL.JoinPath("/api/v2/gpos/", objectID, "/ous")
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
	var finalResponse []OU
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []OU{}, nil
	}
	return finalResponse, nil
}