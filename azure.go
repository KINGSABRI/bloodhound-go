package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetAzureEntity fetches a generic Azure entity by its Object ID.
func (c *Client) GetAzureEntity(objectID string) (json.RawMessage, error) {
	url := c.baseURL.JoinPath("/api/v2/azure/entities/", objectID)
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
		return nil, fmt.Errorf("get azure entity failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get azure entity response: %w", err)
	}
	return response.Data, nil
}

// GetAzureUser fetches a single Azure user by its Object ID.
func (c *Client) GetAzureUser(objectID string) (*AzureUser, *BaseAzureEntity, error) {
	url := c.baseURL.JoinPath("/api/v2/azure/entities/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get azure user failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props      AzureUser       `json:"props"`
			BaseEntity BaseAzureEntity `json:"base"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode get azure user response: %w", err)
	}
	response.Data.Props.ObjectType = "AZUser"
	return &response.Data.Props, &response.Data.BaseEntity, nil
}

// GetAzureUserByName fetches a single Azure user by their User Principal Name.
func (c *Client) GetAzureUserByName(userName string) (*AzureUser, *BaseAzureEntity, error) {
	searchResponse, err := c.Search(userName, "AZUser")
	if err != nil {
		return nil, nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, userName) {
			fullUser, baseEntity, err := c.GetAzureUser(result.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			return fullUser, baseEntity, nil
		}
	}

		return nil, nil, fmt.Errorf("azure user not found: %s", userName)
}

// GetAzureGroup fetches a single Azure group by its Object ID.
func (c *Client) GetAzureGroup(objectID string) (*AzureGroup, *BaseAzureEntity, error) {
	url := c.baseURL.JoinPath("/api/v2/azure/entities/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get azure group failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props      AzureGroup      `json:"props"`
			BaseEntity BaseAzureEntity `json:"base"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode get azure group response: %w", err)
	}
	response.Data.Props.ObjectType = "AZGroup"
	return &response.Data.Props, &response.Data.BaseEntity, nil
}

// GetAzureGroupByName fetches a single Azure group by its name.
func (c *Client) GetAzureGroupByName(groupName string) (*AzureGroup, *BaseAzureEntity, error) {
	searchResponse, err := c.Search(groupName, "AZGroup")
	if err != nil {
		return nil, nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, groupName) {
			fullGroup, baseEntity, err := c.GetAzureGroup(result.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			return fullGroup, baseEntity, nil
		}
	}

	return nil, nil, fmt.Errorf("azure group not found: %s", groupName)
}

// GetAzureVM fetches a single Azure VM by its Object ID.
func (c *Client) GetAzureVM(objectID string) (*AzureVM, *BaseAzureEntity, error) {
	url := c.baseURL.JoinPath("/api/v2/azure/entities/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get azure vm failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props      AzureVM         `json:"props"`
			BaseEntity BaseAzureEntity `json:"base"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, nil, fmt.Errorf("failed to decode get azure vm response: %w", err)
	}
	response.Data.Props.ObjectType = "AZVM"
	return &response.Data.Props, &response.Data.BaseEntity, nil
}

// GetAzureVMByName fetches a single Azure VM by its name.
func (c *Client) GetAzureVMByName(vmName string) (*AzureVM, *BaseAzureEntity, error) {
	searchResponse, err := c.Search(vmName, "AZVM")
	if err != nil {
		return nil, nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, vmName) {
			fullVM, baseEntity, err := c.GetAzureVM(result.ObjectID)
			if err != nil {
				return nil, nil, err
			}
			return fullVM, baseEntity, nil
		}
	}

	return nil, nil, fmt.Errorf("azure vm not found: %s", vmName)
}
