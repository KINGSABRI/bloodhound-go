package bloodhound

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the main BloodHound API client.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	token      string
}

// NewClient creates and returns a new BloodHound API client.
func NewClient(baseURL string) (*Client, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: parsedBaseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 120,
		},
	}, nil
}

// SetHTTPClient allows for overriding the default http.Client.
// This is useful for setting custom transports for logging or testing.
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// setAuthToken sets the session token for the client.
func (c *Client) setAuthToken(token string) {
	c.token = token
}

// SetToken provides a public method to configure the client with a session token.
func (c *Client) SetToken(token string) {
	c.setAuthToken(token)
}

// GetToken returns the session token currently used by the client.
func (c *Client) GetToken() string {
	return c.token
}

// newAuthenticatedRequest creates a new HTTP request with authentication headers.
func (c *Client) newAuthenticatedRequest(method, url string, body io.Reader) (*http.Request, error) {
	if c.token == "" {
		return nil, fmt.Errorf("authentication token is not set")
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "bloodhunter")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes an HTTP request and returns the response.
func (c *Client) do(req *http.Request, sanitizedBody []byte) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read error response body: %w", err)
		}

		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			if len(errorResponse.Errors) > 0 {
				return nil, fmt.Errorf("API error: %s", errorResponse.Errors[0].Message)
			}
		}
		// Fallback to a generic error if parsing fails or there are no specific messages
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Decompress gzipped responses
	if resp.StatusCode != http.StatusNoContent && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader for response: %w", err)
		}
		resp.Body = io.NopCloser(reader)
	}

	return resp, nil
}