package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetDomain fetches a single domain by its Object ID (SID).
func (c *Client) GetDomain(objectID string) (*Domain, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID)
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
		return nil, fmt.Errorf("get domain failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Props Domain `json:"props"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get domain response: %w", err)
	}
	domain := response.Data.Props
	domain.ObjectType = "Domain"
	return &domain, nil
}

// GetDomainByName fetches a single domain by its name.
func (c *Client) GetDomainByName(domainName string) (*Domain, error) {
	searchResponse, err := c.Search(domainName, "Domain")
	if err != nil {
		return nil, err
	}

	for _, result := range searchResponse.Data {
		if strings.EqualFold(result.Name, domainName) {
			return c.GetDomain(result.ObjectID)
		}
	}

	return nil, fmt.Errorf("domain not found: %s", domainName)
}

// GetDomainUsers fetches the users in a given domain.
func (c *Client) GetDomainUsers(objectID string) ([]ADUser, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/users")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []ADUser
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []ADUser{}, nil
	}
	return finalResponse, nil
}

// GetDomainComputers fetches the computers in a given domain.
func (c *Client) GetDomainComputers(objectID string) ([]Computer, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/computers")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Computer
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Computer{}, nil
	}
	return finalResponse, nil
}

// GetDomainForeignUsers fetches the foreign users in a given domain.
func (c *Client) GetDomainForeignUsers(objectID string) ([]ForeignPrincipal, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/foreign-users")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse ForeignPrincipalsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []ForeignPrincipal
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []ForeignPrincipal{}, nil
	}
	return finalResponse, nil
}

// GetDomainInboundTrusts fetches the inbound trusts for a given domain.
func (c *Client) GetDomainInboundTrusts(objectID string) ([]DomainTrust, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/inbound-trusts")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse DomainTrustsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []DomainTrust
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []DomainTrust{}, nil
	}
	return finalResponse, nil
}

// GetDomainOutboundTrusts fetches the outbound trusts for a given domain.
func (c *Client) GetDomainOutboundTrusts(objectID string) ([]DomainTrust, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/outbound-trusts")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse DomainTrustsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []DomainTrust
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []DomainTrust{}, nil
	}
	return finalResponse, nil
}

// GetDomainControllers fetches the controllers of a given domain.
func (c *Client) GetDomainControllers(objectID string) ([]Controller, error) {
	url := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/controllers")
	req, err := c.newAuthenticatedRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawResponse ControllersResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var finalResponse []Controller
	if err := json.Unmarshal(rawResponse.Data, &finalResponse); err != nil {
		return []Controller{}, nil
	}
	return finalResponse, nil
}