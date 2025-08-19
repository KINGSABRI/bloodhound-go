package bloodhound

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the main BloodHound API client.
type Client struct {
	baseURL     *url.URL
	httpClient  *http.Client
	token       string
	debugWriter io.Writer
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
			Timeout: time.Second * 120, // Increased default timeout
		},
	}, nil
}

// SetDebugWriter sets an io.Writer for the client to write debug information to.
func (c *Client) SetDebugWriter(writer io.Writer) {
	c.debugWriter = writer
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

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip") // We accept Gzip responses
	req.Header.Set("Referer", c.baseURL.String())
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Set Content-Type for POST/PUT/PATCH, but do not set Content-Encoding
	if body != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes an HTTP request, handling debug logging and response decompression.
func (c *Client) do(req *http.Request, sanitizedBody []byte) (*http.Response, error) {
	var bodyToLog []byte
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for the actual request
		if sanitizedBody != nil {
			bodyToLog = sanitizedBody
		} else {
			bodyToLog = bodyBytes
		}
	}

	c.logRequest(req, bodyToLog)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	c.logResponse(resp)

	// Manually decompress the body if it was gzipped, because our logResponse function consumed the original body.
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read gzipped response body: %w", err)
		}
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader for response: %w", err)
		}
		resp.Body = io.NopCloser(reader)
	}

	return resp, nil
}