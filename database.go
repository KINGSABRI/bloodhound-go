package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DeleteDatabase deletes all data from the BloodHound database.
func (c *Client) DeleteDatabase() error {
	deleteURL := c.baseURL.JoinPath("/api/v2/clear-database")

	requestBody := DeleteDatabaseRequest{
		DeleteCollectedGraphData: true,
		DeleteFileIngestHistory:  true,
		DeleteDataQualityHistory: true,
	}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal delete database request: %w", err)
	}

	req, err := c.newAuthenticatedRequest(http.MethodPost, deleteURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create delete database request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute delete database request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete database failed with status code: %d", resp.StatusCode)
	}

	return nil
}
