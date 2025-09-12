package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CypherQuery represents a Cypher query.
type CypherQuery struct {
	Query             string `json:"query"`
	IncludeProperties bool   `json:"include_properties"`
}

// SavedQuery represents a saved Cypher query.
type SavedQuery struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Query       string `json:"query"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ListSavedQueries lists all saved Cypher queries.
func (c *Client) ListSavedQueries() ([]SavedQuery, error) {
	var response struct {
		Data []SavedQuery `json:"data"`
	}
	apiUrl := c.baseURL.JoinPath("/api/v2/saved-queries")
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

var ErrDuplicateQueryName = fmt.Errorf("duplicate name for saved query")

// CreateSavedQuery creates a new saved Cypher query.
func (c *Client) CreateSavedQuery(name, query, description string, public bool) (*SavedQuery, error) {
	var savedQuery SavedQuery
	apiUrl := c.baseURL.JoinPath("/api/v2/saved-queries")
	body, err := json.Marshal(SavedQuery{Name: name, Query: query, Description: description, Public: public})
	if err != nil {
		return nil, err
	}
	req, err := c.newAuthenticatedRequest(http.MethodPost, apiUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			for _, e := range errorResponse.Errors {
				if e.Message == "duplicate name for saved query: please choose a different name" {
					return nil, ErrDuplicateQueryName
				}
			}
		}
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create saved query with status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&savedQuery); err != nil {
		return nil, err
	}
	return &savedQuery, nil
}

// UpdateSavedQuery updates a saved Cypher query.
func (c *Client) UpdateSavedQuery(id int, name, query, description string) error {
	apiUrl := c.baseURL.JoinPath(fmt.Sprintf("/api/v2/saved-queries/%d", id))
	body, err := json.Marshal(SavedQuery{Name: name, Query: query, Description: description})
	if err != nil {
		return err
	}
	req, err := c.newAuthenticatedRequest(http.MethodPut, apiUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// DeleteSavedQuery deletes a saved Cypher query.
func (c *Client) DeleteSavedQuery(id int) error {
	apiUrl := c.baseURL.JoinPath(fmt.Sprintf("/api/v2/saved-queries/%d", id))
	req, err := c.newAuthenticatedRequest(http.MethodDelete, apiUrl.String(), nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// ShareSavedQuery shares a saved Cypher query.
func (c *Client) ShareSavedQuery(id int, public bool, userSIDs []string) error {
	apiUrl := c.baseURL.JoinPath(fmt.Sprintf("/api/v2/saved-queries/%d/shares", id))
	body, err := json.Marshal(map[string]interface{}{"public": public, "user_sids": userSIDs})
	if err != nil {
		return err
	}
	req, err := c.newAuthenticatedRequest(http.MethodPost, apiUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// RevokeSavedQuery revokes a saved Cypher query.
func (c *Client) RevokeSavedQuery(id int, userSIDs []string) error {
	apiUrl := c.baseURL.JoinPath(fmt.Sprintf("/api/v2/saved-queries/%d/shares", id))
	body, err := json.Marshal(map[string]interface{}{"user_sids": userSIDs})
	if err != nil {
		return err
	}
	req, err := c.newAuthenticatedRequest(http.MethodDelete, apiUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// RunCypherQuery runs a Cypher query.
func (c *Client) RunCypherQuery(query string) (json.RawMessage, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/graphs/cypher")
	body, err := json.Marshal(CypherQuery{Query: query, IncludeProperties: true})
	if err != nil {
		return nil, err
	}
	req, err := c.newAuthenticatedRequest(http.MethodPost, apiUrl.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			return nil, fmt.Errorf("cypher query failed: %s", errorResponse.Errors[0].Message)
		}
		return nil, fmt.Errorf("cypher query failed with status code: %d", resp.StatusCode)
	}

	if resp.ContentLength == 0 {
		return []byte(`{"data": {}}`), nil
	}

	var response struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode cypher response: %w", err)
	}
	return response.Data, nil
}
