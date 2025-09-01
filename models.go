package bloodhound

import (
	"encoding/json"
	"strconv"
	"time"
)

// JsonTime is a custom time type to handle Unix timestamps from the API.
type JsonTime time.Time

func (jt *JsonTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == "" || s == "0" {
		*jt = JsonTime(time.Time{})
		return nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		*jt = JsonTime(time.Unix(i, 0))
		return nil
	}
	t, err := time.Parse(`"`+time.RFC3339+`"`, s)
	if err != nil {
		*jt = JsonTime(time.Time{})
		return nil
	}
	*jt = JsonTime(t)
	return nil
}

func (jt JsonTime) String() string {
	t := time.Time(jt)
	if t.IsZero() {
		return "<not set>"
	}
	return t.Format("2006-01-02 15:04:05")
}

// LoginRequest represents the JSON body for a login request.
type LoginRequest struct {
	LoginMethod string `json:"login_method"`
	Username    string `json:"username"`
	Secret      string `json:"secret"`
}

// SessionResponse represents the JSON response from a successful login.
type SessionResponse struct {
	Data struct {
		SessionToken string `json:"session_token"`
		UserID       string `json:"user_id"`
		UserDN       string `json:"user_dn"`
	} `json:"data"`
}

// User represents a BloodHound application user object.
type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	IsAdmin   bool     `json:"is_admin"`
	UserDN    string   `json:"user_dn"`
	CreatedAt JsonTime `json:"created_at"`
	UpdatedAt JsonTime `json:"updated_at"`
}

// UsersResponse wraps a list of users, as returned by the API.
type UsersResponse struct {
	Data []User `json:"data"`
}

// CreateUserRequest is the payload for creating a new user.
type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"secret"`
	IsAdmin   bool   `json:"is_admin"`
}

// UpdateUserRequest is the payload for updating an existing user.
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Email     *string `json:"email,omitempty"`
	IsAdmin   *bool   `json:"is_admin,omitempty"`
}

// CypherRequest is the payload for executing a Cypher query.
type CypherRequest struct {
	Query             string `json:"query"`
	IncludeProperties bool   `json:"include_properties"`
}

// CypherResponse wraps the generic response from a Cypher query.
type CypherResponse struct {
	Data json.RawMessage `json:"data"`
}

// CypherResponseData represents the complex data structure returned by a Cypher query.
type CypherResponseData struct {
	Records []map[string]interface{} `json:"records"`
	Nodes   map[string]GraphNodeProperties `json:"nodes"`
	Edges   []GraphEdge                    `json:"edges"`
}

// BaseEntity represents the common properties for all AD objects.
type BaseEntity struct {
	ObjectID          string   `json:"objectid"`
	Name              string   `json:"name"`
	DistinguishedName string   `json:"distinguishedname"`
	SystemTags        string   `json:"system_tags"`
	WhenCreated       JsonTime `json:"whencreated"`
	ObjectType        string   `json:"type"`
}

// Computer represents a BloodHound computer object.
type Computer struct {
	BaseEntity
	OperatingSystem string   `json:"operatingsystem"`
	Enabled         bool     `json:"enabled"`
	HasLAPS         bool     `json:"haslaps"`
	LastSeen        JsonTime `json:"lastseen"`
	// Relationship Counts
	AdminRights    int `json:"adminRights"`
	Controllables  int `json:"controllables"`
	Controllers    int `json:"controllers"`
	DCOMRights     int `json:"dcomRights"`
	PSRemoteRights int `json:"psRemoteRights"`
	RDPRights      int `json:"rdpRights"`
	Sessions       int `json:"sessions"`
}

// ADUser represents a BloodHound AD User object.
type ADUser struct {
	BaseEntity
	Description        string   `json:"description"`
	Enabled            bool     `json:"enabled"`
	HasSIDHistory      bool     `json:"hassidhistory"`
	IsAdmin            bool     `json:"admincount"`
	DontReqPreAuth     bool     `json:"dontreqpreauth"`
	LastLogonTimestamp JsonTime `json:"lastlogontimestamp"`
	// Relationship Counts
	AdminRights           int `json:"adminRights"`
	ConstrainedDelegation int `json:"constrainedDelegation"`
	Controllables         int `json:"controllables"`
	Controllers           int `json:"controllers"`
	DCOMRights            int `json:"dcomRights"`
	GroupMembership       int `json:"groupMembership"`
	PSRemoteRights        int `json:"psRemoteRights"`
	RDPRights             int `json:"rdpRights"`
	Sessions              int `json:"sessions"`
}

// Group represents a BloodHound Group object.
type Group struct {
	BaseEntity
	IsAdmin      bool   `json:"admincount"`
	Description  string `json:"description"`
	SamAccountName string `json:"samaccountname"`
	// Relationship Counts
	AdminRights    int `json:"adminRights"`
	Controllables  int `json:"controllables"`
	Controllers    int `json:"controllers"`
	DCOMRights     int `json:"dcomRights"`
	Members        int `json:"members"`
	PSRemoteRights int `json:"psRemoteRights"`
	RDPRights      int `json:"rdpRights"`
	Sessions       int `json:"sessions"`
}

// Domain represents a BloodHound Domain object.
type Domain struct {
	BaseEntity
}

// GPO represents a BloodHound GPO object.
type GPO struct {
	BaseEntity
}

// OU represents a BloodHound OU object.
type OU struct {
	BaseEntity
}

// BaseAzureEntity represents the common properties for all Azure objects.
type BaseAzureEntity struct {
	ObjectID   string `json:"objectid"`
	Name       string `json:"name"`
	SystemTags string `json:"system_tags"`
	ObjectType string `json:"type"`
}

// AzureUser represents a BloodHound Azure User object.
type AzureUser struct {
	BaseAzureEntity
	UserPrincipalName string `json:"userprincipalname"`
	Enabled           bool   `json:"enabled"`
}

// AzureGroup represents a BloodHound Azure Group object.
type AzureGroup struct {
	BaseAzureEntity
}

// AzureVM represents a BloodHound Azure VM object.
type AzureVM struct {
	BaseAzureEntity
	OperatingSystem string `json:"operatingsystem"`
}

// EntityAdmin represents a principal with admin rights to another entity.
type EntityAdmin struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// EntityAdminsResponse wraps a list of entity admins.
type EntityAdminsResponse struct {
	Data json.RawMessage `json:"data"`
}

// Session represents a user session on a computer.
type Session struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// SessionsResponse wraps a list of sessions.
type SessionsResponse struct {
	Data json.RawMessage `json:"data"`
}

// Privilege represents a principal with a specific privilege on an entity.
type Privilege struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// PrivilegesResponse wraps a list of privileges.
type PrivilegesResponse struct {
	Data json.RawMessage `json:"data"`
}

// ConstrainedDelegation represents a constrained delegation privilege.
type ConstrainedDelegation struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ConstrainedDelegationsResponse wraps a list of constrained delegations.
type ConstrainedDelegationsResponse struct {
	Data json.RawMessage `json:"data"`
}

// GroupMembership represents a group that a principal is a member of.
type GroupMembership struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// GroupMembershipsResponse wraps a list of group memberships.
type GroupMembershipsResponse struct {
	Data json.RawMessage `json:"data"`
}

// FileUploadJob represents a file upload job.
type FileUploadJob struct {
	ID int `json:"id"`
}

// FileUploadResponse wraps the response from a file upload.
type FileUploadResponse struct {
	Data FileUploadJob `json:"data"`
}

// DomainTrust represents a trust relationship between domains.
type DomainTrust struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// DomainTrustsResponse wraps a list of domain trusts.
type DomainTrustsResponse struct {
	Data json.RawMessage `json:"data"`
}

// Controller represents a principal that controls another entity.
type Controller struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ControllersResponse wraps a list of controllers.
type ControllersResponse struct {
	Data json.RawMessage `json:"data"`
}

// Controllable represents an entity that is controlled by another principal.
type Controllable struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ControllablesResponse wraps a list of controllables.
type ControllablesResponse struct {
	Data json.RawMessage `json:"data"`
}

// SearchResult represents a single result from the search endpoint.
type SearchResult struct {
	ObjectID   string `json:"objectid"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	SystemTags string `json:"system_tags"`
}

// SearchResponse wraps the response from the search endpoint.
type SearchResponse struct {
	Data []SearchResult `json:"data"`
}

// AttackPath represents a single attack path.
type AttackPath struct {
	ObjectIDs   []string `json:"object_ids"`
	NodeCounts  map[string]int `json:"node_counts"`
	Severity    float64  `json:"severity"`
	PathFinding json.RawMessage `json:"path_finding"`
}

// AttackPathsResponse wraps a list of attack paths.
type AttackPathsResponse struct {
	Data json.RawMessage `json:"data"`
}

// AttackPathFinding represents a single attack path finding.
type AttackPathFinding struct {
	Finding json.RawMessage `json:"finding"`
}

// AttackPathFindingsResponse wraps a list of attack path findings.
type AttackPathFindingsResponse struct {
	Data json.RawMessage `json:"data"`
}

// ForeignPrincipal represents a user or group from a foreign domain.
type ForeignPrincipal struct {
	Name       string `json:"name"`
	ObjectID   string `json:"object_id"`
	ObjectType string `json:"object_type"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ForeignPrincipalsResponse wraps a list of foreign principals.
type ForeignPrincipalsResponse struct {
	Data json.RawMessage `json:"data"`
}

// GraphNodeProperties represents the properties of a node in a graph response.
type GraphNodeProperties struct {
	Name       string `json:"label"`
	Kind       string `json:"kind"`
	ObjectID   string `json:"objectId"`
	SystemTags string `json:"system_tags"`
}

// GraphEdge represents an edge in a graph response.
type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
	Kind   string `json:"kind"`
}

// ShortestPathData represents the nodes and edges of the shortest path graph.
type ShortestPathData struct {
	Nodes map[string]GraphNodeProperties `json:"nodes"`
	Edges []GraphEdge                    `json:"edges"`
}

// ShortestPathResponse wraps the response from a shortest path query.
type ShortestPathResponse struct {
	Data ShortestPathData `json:"data"`
}

// AssetGroup represents a BloodHound asset group.
type AssetGroup struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Tag         string     `json:"tag"`
	SystemGroup bool       `json:"system_group"`
	Selectors   []Selector `json:"Selectors"`
}

// Selector represents a selector for an asset group.
type Selector struct {
	AssetGroupID int    `json:"asset_group_id"`
	Name         string `json:"name"`
	Selector     string `json:"selector"`
}

// AssetGroupsResponse wraps a list of asset groups.
type AssetGroupsResponse struct {
	Data struct {
		AssetGroups []AssetGroup `json:"asset_groups"`
	} `json:"data"`
}

// OwnershipUpdate represents a single update to the Owned asset group.
type OwnershipUpdate struct {
	SelectorName string `json:"selector_name"`
	SID          string `json:"sid"`
	Action       string `json:"action"`
}