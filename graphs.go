package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetShortestPath calculates the shortest path between two nodes in the graph,
// optionally filtering by a list of relationship kinds.
func (c *Client) GetShortestPath(startNodeID, endNodeID string, relationshipKinds []string) (*ShortestPathResponse, error) {
	shortestPathURL := c.baseURL.JoinPath("/api/v2/graphs/shortest-path")
	params := url.Values{}
	params.Add("start_node", startNodeID)
	params.Add("end_node", endNodeID)
	if len(relationshipKinds) > 0 {
		params.Add("relationship_kinds", strings.Join(relationshipKinds, ","))
	}
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
		return nil, fmt.Errorf("shortest path query failed with status code: %d", resp.StatusCode)
	}

	var shortestPathResponse ShortestPathResponse
	if err := json.NewDecoder(resp.Body).Decode(&shortestPathResponse); err != nil {
		return nil, fmt.Errorf("failed to decode shortest path response: %w", err)
	}

	return &shortestPathResponse, nil
}