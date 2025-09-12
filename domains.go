package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetDomain fetches a single domain by its Object ID (SID).
func (c *Client) GetDomain(objectID string) (*Domain, error) {
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
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
			Props                 Domain `json:"props"`
			Users                       int      `json:"users"`
			Computers                   int      `json:"computers"`
			Controllers                 int      `json:"controllers"`
			DCSyncers                   int      `json:"dcsyncers"`
			ForeignAdmins               int      `json:"foreignAdmins"`
			ForeignGPOControllers       int      `json:"foreignGPOControllers"`
			ForeignGroups               int      `json:"foreignGroups"`
			ForeignUsers                int      `json:"foreignUsers"`
			GPOs                        int      `json:"gpos"`
			Groups                      int      `json:"groups"`
			InboundTrusts               int      `json:"inboundTrusts"`
			LinkedGPOs                  int      `json:"linkedgpos"`
			OUs                         int      `json:"ous"`
			OutboundTrusts              int      `json:"outboundTrusts"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode get domain response: %w", err)
	}
	domain := response.Data.Props
	domain.ObjectType = "Domain"
	domain.Users = response.Data.Users
	domain.Computers = response.Data.Computers
	domain.Controllers = response.Data.Controllers
	domain.DCSyncers = response.Data.DCSyncers
	domain.ForeignAdmins = response.Data.ForeignAdmins
	domain.ForeignGPOControllers = response.Data.ForeignGPOControllers
	domain.ForeignGroups = response.Data.ForeignGroups
	domain.ForeignUsers = response.Data.ForeignUsers
	domain.GPOs = response.Data.GPOs
	domain.Groups = response.Data.Groups
	domain.InboundTrusts = response.Data.InboundTrusts
	domain.LinkedGPOs = response.Data.LinkedGPOs
	domain.OUs = response.Data.OUs
	domain.OutboundTrusts = response.Data.OutboundTrusts
	return &domain, nil
}

// GetDomainByName fetches a single domain by its name.
func (c *Client) GetDomainByName(domainName string) (*Domain, error) {
	searchResponse, err := c.Search(domainName, "Domain", 0)
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
func (c *Client) GetDomainUsers(objectID string, limit int) (UsersResponse, error) {
	var rawResponse UsersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/users")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainComputers fetches the computers in a given domain.
func (c *Client) GetDomainComputers(objectID string, limit int) (ComputersResponse, error) {
	var rawResponse ComputersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/computers")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainForeignUsers fetches the foreign users in a given domain.
func (c *Client) GetDomainForeignUsers(objectID string, limit int) (ForeignPrincipalsResponse, error) {
	var rawResponse ForeignPrincipalsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/foreign-users")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainInboundTrusts fetches the inbound trusts for a given domain.
func (c *Client) GetDomainInboundTrusts(objectID string, limit int) (DomainTrustsResponse, error) {
	var rawResponse DomainTrustsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/inbound-trusts")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainOutboundTrusts fetches the outbound trusts for a given domain.
func (c *Client) GetDomainOutboundTrusts(objectID string, limit int) (DomainTrustsResponse, error) {
	var rawResponse DomainTrustsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/outbound-trusts")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainOUs fetches the OUs in a given domain.
func (c *Client) GetDomainOUs(objectID string, limit int) (OUsResponse, error) {
	var rawResponse OUsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/ous")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainLinkedGPOs fetches the linked GPOs in a given domain.
func (c *Client) GetDomainLinkedGPOs(objectID string, limit int) (GPOsResponse, error) {
	var rawResponse GPOsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/linked-gpos")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainGroups fetches the groups in a given domain.
func (c *Client) GetDomainGroups(objectID string, limit int) (GroupsResponse, error) {
	var rawResponse GroupsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/groups")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainGPOs fetches the GPOs in a given domain.
func (c *Client) GetDomainGPOs(objectID string, limit int) (GPOsResponse, error) {
	var rawResponse GPOsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/gpos")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainForeignGroups fetches the foreign groups in a given domain.
func (c *Client) GetDomainForeignGroups(objectID string, limit int) (ForeignPrincipalsResponse, error) {
	var rawResponse ForeignPrincipalsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/foreign-groups")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainForeignGPOControllers fetches the foreign GPO controllers in a given domain.
func (c *Client) GetDomainForeignGPOControllers(objectID string, limit int) (ForeignPrincipalsResponse, error) {
	var rawResponse ForeignPrincipalsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/foreign-gpo-controllers")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainForeignAdmins fetches the foreign admins in a given domain.
func (c *Client) GetDomainForeignAdmins(objectID string, limit int) (ForeignPrincipalsResponse, error) {
	var rawResponse ForeignPrincipalsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/foreign-admins")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainDCSyncers fetches the principals with DCSync rights to a given domain.
func (c *Client) GetDomainDCSyncers(objectID string, limit int) (ForeignPrincipalsResponse, error) {
	var rawResponse ForeignPrincipalsResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/dc-syncers")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// GetDomainControllers fetches the controllers of a given domain.
func (c *Client) GetDomainControllers(objectID string, limit int) (ControllersResponse, error) {
	var rawResponse ControllersResponse
	apiUrl := c.baseURL.JoinPath("/api/v2/domains/", objectID, "/controllers")
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", limit))
	}
	apiUrl.RawQuery = params.Encode()
	req, err := c.newAuthenticatedRequest(http.MethodGet, apiUrl.String(), nil)
	if err != nil {
		return rawResponse, err
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return rawResponse, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return rawResponse, err
	}
	return rawResponse, nil
}

// ListDomains fetches all domains.
func (c *Client) ListDomains() ([]AvailableDomain, error) {
	var response struct {
		Data []AvailableDomain `json:"data"`
	}
	apiUrl := c.baseURL.JoinPath("/api/v2/available-domains")
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