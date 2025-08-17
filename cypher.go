package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RunCypherQuery executes a raw Cypher query against the BloodHound graph.
func (c *Client) RunCypherQuery(query string) (*CypherResponse, error) {
	cypherURL := c.baseURL.JoinPath("/api/v2/graphs/cypher")

	requestBody := CypherRequest{
		Query:             query,
		IncludeProperties: true,
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cypher request: %w", err)
	}

	req, err := c.newAuthenticatedRequest(http.MethodPost, cypherURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create cypher request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute cypher request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cypher query failed with status code: %d", resp.StatusCode)
	}

	var cypherResponse CypherResponse
	if err := json.NewDecoder(resp.Body).Decode(&cypherResponse); err != nil {
		return nil, fmt.Errorf("failed to decode cypher response: %w", err)
	}

	return &cypherResponse, nil
}