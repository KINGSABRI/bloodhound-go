package bloodhound

import (
	"encoding/json"
	"net/http"
)

// ListAttackPaths fetches the list of available attack paths.
func (c *Client) ListAttackPaths() ([]AttackPath, error) {
	url := c.baseURL.JoinPath("/api/v2/attack-paths")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response AttackPathsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

// ListAttackPathFindings fetches the list of attack path findings.
func (c *Client) ListAttackPathFindings() ([]AttackPathFinding, error) {
	url := c.baseURL.JoinPath("/api/v2/attack-paths/findings")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response AttackPathFindingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data, nil
}
