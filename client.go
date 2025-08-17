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
			Timeout: time.Second * 30,
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

	// Set common browser headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Referer", c.baseURL.String())

	// Set authentication and content type
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes an HTTP request, handling Gzipping, debug logging, and returns the response.
func (c *Client) do(req *http.Request, sanitizedBody []byte) (*http.Response, error) {
	// Read the original body for potential Gzipping and logging
	var originalBody []byte
	if req.Body != nil {
		originalBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(originalBody)) // Restore body for the actual request
	}

	// Gzip the body for POST, PATCH, PUT requests
	if req.Method == http.MethodPost || req.Method == http.MethodPatch || req.Method == http.MethodPut {
		if len(originalBody) > 0 && req.Header.Get("Content-Type") == "application/json" {
			var compressedBody bytes.Buffer
			gz := gzip.NewWriter(&compressedBody)
			if _, err := gz.Write(originalBody); err != nil {
				return nil, fmt.Errorf("failed to gzip request body: %w", err)
			}
			if err := gz.Close(); err != nil {
				return nil, fmt.Errorf("failed to close gzip writer: %w", err)
			}
			req.Body = io.NopCloser(&compressedBody)
			req.ContentLength = int64(compressedBody.Len())
			req.Header.Set("Content-Encoding", "gzip")
		}
	}

	// Use the sanitized body for logging if provided, otherwise use the original.
	var logBody []byte
	if sanitizedBody != nil {
		logBody = sanitizedBody
	} else {
		logBody = originalBody
	}

	c.logRequest(req, logBody)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// The http.Client transport automatically decompresses the body if Content-Encoding is gzip
	// but we need to log the raw response first.
	c.logResponse(resp)

	// Now, we need to manually decompress the body for the application if it was gzipped,
	// because our logResponse function consumed the original (potentially compressed) body.
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