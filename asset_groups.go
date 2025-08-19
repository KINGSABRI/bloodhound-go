package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// getOwnedAssetGroupID fetches all asset groups and returns the ID of the "Owned" group.
func (c *Client) getOwnedAssetGroupID() (int, error) {
	assetGroupsURL := c.baseURL.JoinPath("/api/v2/asset-groups")
	req, err := c.newAuthenticatedRequest(http.MethodGet, assetGroupsURL.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create get asset groups request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to execute get asset groups request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("get asset groups failed with status code: %d", resp.StatusCode)
	}

	var assetGroupsResponse AssetGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetGroupsResponse); err != nil {
		return 0, fmt.Errorf("failed to decode asset groups response: %w", err)
	}

	for _, group := range assetGroupsResponse.Data.AssetGroups {
		if group.Name == "Owned" {
			return group.ID, nil
		}
	}

	return 0, fmt.Errorf("could not find the 'Owned' asset group")
}

// UpdateOwnedStatus adds or removes objects from the 'Owned' asset group.
func (c *Client) UpdateOwnedStatus(updates []OwnershipUpdate) error {
	ownedGroupID, err := c.getOwnedAssetGroupID()
	if err != nil {
		return err
	}

	updateURL := c.baseURL.JoinPath("/api/v2/asset-groups/", fmt.Sprintf("%d", ownedGroupID), "/selectors")

	payload, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("failed to marshal ownership update payload: %w", err)
	}

	req, err := c.newAuthenticatedRequest(http.MethodPut, updateURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create ownership update request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute ownership update request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("ownership update failed with status code: %d", resp.StatusCode)
	}

	return nil
}
