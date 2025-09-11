package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetShortestPath finds the shortest path between two nodes.
func (c *Client) GetShortestPath(startNode, endNode, relationshipKinds string) (*ShortestPathResponse, error) {
	params := url.Values{}
	params.Add("start_node", startNode)
	params.Add("end_node", endNode)
	params.Add("relationship_kinds", relationshipKinds)

	shortestPathURL := c.baseURL.JoinPath("/api/v2/graphs/shortest-path")
	shortestPathURL.RawQuery = params.Encode()

	req, err := c.newAuthenticatedRequest(http.MethodGet, shortestPathURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create shortest path request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute shortest path request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			var errorResponse ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
				if len(errorResponse.Errors) > 0 && errorResponse.Errors[0].Message == "Path not found" {
					return nil, fmt.Errorf("path not found")
				}
			}
		}
		return nil, fmt.Errorf("shortest path request failed with status code: %d", resp.StatusCode)
	}

	var shortestPathResponse ShortestPathResponse
	if err := json.NewDecoder(resp.Body).Decode(&shortestPathResponse); err != nil {
		return nil, fmt.Errorf("failed to decode shortest path response: %w", err)
	}

	return &shortestPathResponse, nil
}
