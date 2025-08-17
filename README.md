#  bloodhound-go

[![Go Version](https://img.shields.io/badge/go-1.22+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, type-safe, and production-ready Go client for the official [BloodHound API](https://bloodhound.specterops.io/reference/).

`bloodhound-go` provides a clean and idiomatic Go interface for interacting with all major aspects of the BloodHound API, making it easy to build custom tools, scripts, and integrations.

---

## Features

- **Full API Coverage:** Implements the vast majority of BloodHound API endpoints.
  - Authentication & User Management
  - AD & Azure Entity Querying (Users, Computers, Groups, GPOs, Domains, OUs)
  - Raw Cypher Query Execution
  - Multi-step Data Ingestion (File Upload)
  - Attack Path Analysis
- **Type-Safe:** All API requests and responses are mapped to Go structs, preventing common errors and improving developer experience.
- **Robust Error Handling:** Provides clear and descriptive errors for failed API calls.
- **Clean Interface:** Designed to be intuitive and easy to integrate into your Go projects.

## Installation

```bash
go get github.com/user/bloodhound-go
```

## Usage Example

Here's a quick example of how to use the client to log in, fetch the current user's information, and log out.

```go
package main

import (
	"fmt"
	"log"

	"github.com/user/bloodhound-go"
)

func main() {
	// --- 1. Initialize the Client ---
	// Replace with your BloodHound instance URL.
	bhClient, err := bloodhound.NewClient("http://localhost:8080")
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// --- 2. Authenticate ---
	// Use a service account or user credentials.
	err = bhClient.Login("username", "password")
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("Successfully authenticated to BloodHound.")

	// --- 3. Make Authenticated API Calls ---
	// Example: Get information about the current user.
	self, err := bhClient.GetSelf()
	if err != nil {
		log.Fatalf("Failed to get self: %v", err)
	}
	fmt.Printf("Current User: %s %s (%s)\n", self.FirstName, self.LastName, self.Email)

	// Example: List all users in BloodHound.
	users, err := bhClient.ListUsers()
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	fmt.Printf("Found %d total users.\n", len(users))


	// --- 4. Log Out ---
	// Invalidate the session token.
	if err := bhClient.Logout(); err != nil {
		log.Fatalf("Logout failed: %v", err)
	}
	fmt.Println("Successfully logged out.")
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```