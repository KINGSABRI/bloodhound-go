package bloodhound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSelf fetches the user object for the currently authenticated user.
func (c *Client) GetSelf() (User, error) {
	var user User
	selfURL := c.baseURL.JoinPath("/api/v2/self")
	req, err := c.newAuthenticatedRequest(http.MethodGet, selfURL.String(), nil)
	if err != nil {
		return user, fmt.Errorf("failed to create self request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return user, fmt.Errorf("failed to execute self request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("get self failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			User   User   `json:"user"`
			UserDN string `json:"user_dn"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return user, fmt.Errorf("failed to decode self response: %w", err)
	}

	// Manually set the UserDN on the returned user object
	user = response.Data.User
	user.UserDN = response.Data.UserDN

	return user, nil
}

// ListUsers fetches all users from the BloodHound instance.
func (c *Client) ListUsers() ([]User, error) {
	usersURL := c.baseURL.JoinPath("/api/v2/users")
	req, err := c.newAuthenticatedRequest(http.MethodGet, usersURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list users request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list users request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list users failed with status code: %d", resp.StatusCode)
	}

	var usersResponse UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
		return nil, fmt.Errorf("failed to decode list users response: %w", err)
	}

	return usersResponse.Data, nil
}

// GetUser fetches a single user by their ID.
func (c *Client) GetUser(userID string) (User, error) {
	var user User
	userURL := c.baseURL.JoinPath("/api/v2/users/", userID)
	req, err := c.newAuthenticatedRequest(http.MethodGet, userURL.String(), nil)
	if err != nil {
		return user, fmt.Errorf("failed to create get user request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return user, fmt.Errorf("failed to execute get user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("get user failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return user, fmt.Errorf("failed to decode get user response: %w", err)
	}

	return response.Data, nil
}

// CreateUser creates a new BloodHound user.
func (c *Client) CreateUser(request CreateUserRequest) (User, error) {
	var user User
	payload, err := json.Marshal(request)
	if err != nil {
		return user, fmt.Errorf("failed to marshal create user request: %w", err)
	}

	usersURL := c.baseURL.JoinPath("/api/v2/users")
	req, err := c.newAuthenticatedRequest(http.MethodPost, usersURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return user, fmt.Errorf("failed to create create user request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return user, fmt.Errorf("failed to execute create user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return user, fmt.Errorf("create user failed with status code: %d", resp.StatusCode)
	}

	var response struct {
		Data User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return user, fmt.Errorf("failed to decode create user response: %w", err)
	}

	return response.Data, nil
}

// UpdateUser updates an existing BloodHound user.
func (c *Client) UpdateUser(userID string, request UpdateUserRequest) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal update user request: %w", err)
	}

	userURL := c.baseURL.JoinPath("/api/v2/users/", userID)
	req, err := c.newAuthenticatedRequest(http.MethodPatch, userURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create update user request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute update user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update user failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// DeleteUser deletes a BloodHound user.
func (c *Client) DeleteUser(userID string) error {
	userURL := c.baseURL.JoinPath("/api/v2/users/", userID)
	req, err := c.newAuthenticatedRequest(http.MethodDelete, userURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete user request: %w", err)
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute delete user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete user failed with status code: %d", resp.StatusCode)
	}

	return nil
}
