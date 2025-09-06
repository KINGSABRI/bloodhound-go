package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Search looks for an object in BloodHound by its name, optionally filtering by type and limit.
func (c *Client) Search(searchTerm, objectType string, limit int) (*SearchResponse, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/search")
	params := url.Values{}
	params.Add("q", searchTerm)
	if objectType != "" {
		params.Add("type", objectType)
	}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()

	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request failed with status code: %d", resp.StatusCode)
	}

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &searchResponse, nil
}
