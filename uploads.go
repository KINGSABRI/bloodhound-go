package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// StartFileUploadJob starts a new file upload job.
func (c *Client) StartFileUploadJob() (*FileUploadJob, error) {
	startURL := c.baseURL.JoinPath("/api/v2/file-upload/start")
	req, err := c.newAuthenticatedRequest(http.MethodPost, startURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create start file upload request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute start file upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("start file upload failed with status code: %d", resp.StatusCode)
	}

	var fileUploadResponse FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileUploadResponse); err != nil {
		return nil, fmt.Errorf("failed to decode start file upload response: %w", err)
	}

	return &fileUploadResponse.Data, nil
}

// UploadFile uploads a file to a file upload job.
func (c *Client) UploadFile(jobID int, content []byte, contentType string) error {
	uploadURL := c.baseURL.JoinPath("/api/v2/file-upload/", strconv.Itoa(jobID))

	req, err := c.newAuthenticatedRequest(http.MethodPost, uploadURL.String(), bytes.NewBuffer(content))
	if err != nil {
		return fmt.Errorf("failed to create upload file request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute upload file request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload file failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// EndFileUploadJob marks a file upload job as complete.
func (c *Client) EndFileUploadJob(jobID int) error {
	endURL := c.baseURL.JoinPath("/api/v2/file-upload/", strconv.Itoa(jobID), "/end")
	req, err := c.newAuthenticatedRequest(http.MethodPost, endURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create end file upload request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute end file upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("end file upload failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListFileUploadJobs lists all file upload jobs.
func (c *Client) ListFileUploadJobs() ([]FileUploadJob, error) {
	listURL := c.baseURL.JoinPath("/api/v2/file-upload")
	req, err := c.newAuthenticatedRequest(http.MethodGet, listURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list file upload jobs request: %w", err)
	}

	var jobsResponse struct {
		Data []FileUploadJob `json:"data"`
	}
	resp, err := c.do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list file upload jobs request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list file upload jobs failed with status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&jobsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode list file upload jobs response: %w", err)
	}

	return jobsResponse.Data, nil
}
