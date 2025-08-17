package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetShortestPath calculates the shortest path between two nodes in the graph.
func (c *Client) GetShortestPath(startNodeID, endNodeID string) (*ShortestPathResponse, error) {
	shortestPathURL := c.baseURL.JoinPath("/api/v2/graphs/shortest-path")

	requestBody := ShortestPathRequest{
		StartNode: startNodeID,
		EndNode:   endNodeID,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal shortest path request: %w", err)
	}

	req, err := c.newAuthenticatedRequest(http.MethodPost, shortestPathURL.String(), bytes.NewBuffer(payload))
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
