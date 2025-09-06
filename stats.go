package bloodhound

import (
	"encoding/json"
	"fmt"
)

// Stat represents a single statistic from the database.
type Stat struct {
	ObjectType string `json:"type"`
	Count      int    `json:"count"`
}

// GetGlobalStats fetches the total count of all object types in the database.
func (c *Client) GetGlobalStats() ([]Stat, error) {
	query := `MATCH (n) RETURN labels(n) as type, count(n) as count`
	resp, err := c.RunCypherQuery(query)
	if err != nil {
		return nil, err
	}

	var stats []Stat
	if err := json.Unmarshal(resp, &stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats response: %w", err)
	}
	return stats, nil
}

// DomainStats represents a collection of detailed statistics for a single domain.
type DomainStats struct {
	TotalUsers         int
	AdminUsers         int
	KerberoastableUsers int
	TotalComputers     int
	DomainControllers  int
	TotalGroups        int
	TotalGPOs          int
	TotalOUs           int
}

// GetDomainStats fetches detailed statistics for a specific domain.
func (c *Client) GetDomainStats(domainName string) (*DomainStats, error) {
	stats := &DomainStats{}
	var err error

	// Helper to run a query and populate a stat field
	runQuery := func(field *int, query string, params map[string]interface{}) error {
		// We don't have a parameterized query function yet, so we'll do simple replacement for now.
		// This is safe because the domain name is the only parameter.
		query = fmt.Sprintf(query, domainName)
		resp, err := c.RunCypherQuery(query)
		if err != nil {
			return err
		}
		var result []struct {
			Count int `json:"count"`
		}
		if err := json.Unmarshal(resp, &result); err != nil {
			return err
		}
		if len(result) > 0 {
			*field = result[0].Count
		}
		return nil
	}

	// Run all the queries
	err = runQuery(&stats.TotalUsers, `MATCH (u:User) WHERE u.domain = "%s" RETURN count(u) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.AdminUsers, `MATCH (u:User) WHERE u.domain = "%s" AND u.admincount = true RETURN count(u) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.KerberoastableUsers, `MATCH (u:User) WHERE u.domain = "%s" AND u.hasspn = true RETURN count(u) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.TotalComputers, `MATCH (c:Computer) WHERE c.domain = "%s" RETURN count(c) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.DomainControllers, `MATCH (c:Computer) WHERE c.domain = "%s" AND c.operatingsystem CONTAINS "Server" RETURN count(c) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.TotalGroups, `MATCH (g:Group) WHERE g.domain = "%s" RETURN count(g) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.TotalGPOs, `MATCH (g:GPO) WHERE g.domain = "%s" RETURN count(g) as count`, nil)
	if err != nil { return nil, err }
	err = runQuery(&stats.TotalOUs, `MATCH (o:OU) WHERE o.domain = "%s" RETURN count(o) as count`, nil)
	if err != nil { return nil, err }

	return stats, nil
}
