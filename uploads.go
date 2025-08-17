package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("start file upload failed with status code: %d", resp.StatusCode)
	}

	var fileUploadResponse FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileUploadResponse); err != nil {
		return nil, fmt.Errorf("failed to decode start file upload response: %w", err)
	}

	return &fileUploadResponse.Data, nil
}

// UploadFile uploads a chunk of a file to a file upload job.
func (c *Client) UploadFile(jobID int, chunk []byte, chunkNumber, totalChunks int) error {
	uploadURL := c.baseURL.JoinPath("/api/v2/file-upload/", strconv.Itoa(jobID))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "chunk.zip")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, bytes.NewReader(chunk)); err != nil {
		return fmt.Errorf("failed to copy chunk to form file: %w", err)
	}

	writer.WriteField("resumableChunkNumber", strconv.Itoa(chunkNumber))
	writer.WriteField("resumableTotalChunks", strconv.Itoa(totalChunks))
	writer.Close()

	req, err := c.newAuthenticatedRequest(http.MethodPost, uploadURL.String(), body)
	if err != nil {
		return fmt.Errorf("failed to create upload file request: %w", err)
	}
	// This is a special case; override the content type.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.do(req, nil) // We don't log the binary chunk body
	if err != nil {
		return fmt.Errorf("failed to execute upload file request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload file failed with status code: %d", resp.StatusCode)
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("end file upload failed with status code: %d", resp.StatusCode)
	}

	return nil
}