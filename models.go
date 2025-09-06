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
	ID            string   `json:"id"`
	Email         string   `json:"email_address"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	PrincipalName string   `json:"principal_name"`
	IsAdmin       bool     `json:"is_admin"`
	UserDN        string   `json:"user_dn"`
	CreatedAt     JsonTime `json:"created_at"`
	UpdatedAt     JsonTime `json:"updated_at"`
}

// UsersResponse wraps a list of users, as returned by the API.
type UsersResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
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

// Computer represents a computer in BloodHound.
type Computer struct {
	Name                    string   `json:"name"`
	ObjectID                string   `json:"objectid"`
	ObjectType              string   `json:"object_type"`
	DistinguishedName       string   `json:"distinguishedname"`
	OperatingSystem         string   `json:"operatingsystem"`
	Enabled                 bool     `json:"enabled"`
	HasLAPS                 bool     `json:"haslaps"`
	WhenCreated             interface{} `json:"whencreated"`
	LastSeen                string   `json:"lastseen"`
	ServicePrincipalNames   []string `json:"serviceprincipalnames"`
	UnconstrainedDelegation bool     `json:"unconstraineddelegation"`
	SupportedEncryptionTypes []string `json:"supportedencryptiontypes"`
	LastLogon               interface{} `json:"lastlogon"`
	IsDC                    bool     `json:"isdc"`
	AdminRights             int      `json:"adminRights"`
	AdminUsers              int      `json:"adminUsers"`
	ConstrainedPrivs        int      `json:"constrainedPrivs"`
	ConstrainedUsers        int      `json:"constrainedUsers"`
	Controllables           int      `json:"controllables"`
	Controllers             int      `json:"controllers"`
	DCOMRights              int      `json:"dcomRights"`
	DCOMUsers               int      `json:"dcomUsers"`
	GPOs                    int      `json:"gpos"`
	GroupMembership         int      `json:"groupMembership"`
	PSRemoteRights          int      `json:"psRemoteRights"`
	PSRemoteUsers           int      `json:"psRemoteUsers"`
	RDPRights               int      `json:"rdpRights"`
	Sessions                int      `json:"sessions"`
	SQLAdminUsers           int      `json:"sqlAdminUsers"`
}

// ADUser represents a BloodHound AD User object.
type ADUser struct {
	BaseEntity
	DisplayName             string   `json:"displayname"`
	Description             string   `json:"description"`
	Email                   string   `json:"email"`
	Domain                  string   `json:"domain"`
	DomainSID               string   `json:"domainsid"`
	Enabled                 bool     `json:"enabled"`
	HasSIDHistory           bool     `json:"hassidhistory"`
	IsAdmin                 bool     `json:"admincount"`
	DontReqPreAuth          bool     `json:"dontreqpreauth"`
	LockedOut               bool     `json:"lockedout"`
	PwdNeverExpires         bool     `json:"pwdneverexpires"`
	UnconstrainedDelegation bool     `json:"unconstraineddelegation"`
	HasSPN                  bool     `json:"hasspn"`
	ServicePrincipalNames   []string `json:"serviceprincipalnames"`
	LastLogonTimestamp      JsonTime `json:"lastlogontimestamp"`
	LastSeen                JsonTime `json:"lastseen"`
	PwdLastSet              JsonTime `json:"pwdlastset"`
	// Relationship Counts
	AdminRights           int `json:"adminRights"`
	ConstrainedDelegation int `json:"constrainedDelegation"`
	Controllables         int `json:"controllables"`
	Controllers           int `json:"controllers"`
	DCOMRights            int `json:"dcomRights"`
	GroupMembership       int `json:"groupMembership"`
	GPOs                  int `json:"gpos"`
	PSRemoteRights        int `json:"psRemoteRights"`
	RDPRights             int `json:"rdpRights"`
	Sessions              int `json:"sessions"`
	SQLAdminRights        int `json:"sqlAdminRights"`
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

// Domain represents a domain in BloodHound.
type Domain struct {
	Name                        string   `json:"name"`
	ObjectID                    string   `json:"objectid"`
	ObjectType                  string   `json:"object_type"`
	DistinguishedName           string   `json:"distinguishedname"`
	WhenCreated                 interface{} `json:"whencreated"`
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
	FSMORoleOwner               string   `json:"fsmoroleowner"`
	DC                          string   `json:"dc"`
	FunctionalLevel             string   `json:"functionallevel"`
	MachineAccountQuota         int      `json:"ms-ds-machineaccountquota"`
	AllUsersTrustQuota          int      `json:"msds-alluserstrustquota"`
	IsACLProtected              bool     `json:"isaclprotected"`
	IsCriticalSystemObject      bool     `json:"iscriticalsystemobject"`
	LastCollected               string   `json:"lastcollected"`
	LastSeen                    string   `json:"lastseen"`
	LockoutDuration             interface{} `json:"lockoutduration"`
	LockoutObservationWindow    interface{} `json:"lockoutobservationwindow"`
	LockoutThreshold            int      `json:"lockoutthreshold"`
	MasteredBy                  []string `json:"masteredby"`
	MSDSMasteredBy              []string `json:"msds-masteredby"`
	WellKnownObjects            []string `json:"wellknownobjects"`
}

// GPO represents a BloodHound GPO object.
type GPO struct {
	Name              string      `json:"name"`
	ObjectID          string      `json:"objectid"`
	ObjectType        string      `json:"object_type"`
	DistinguishedName string      `json:"distinguishedname"`
	WhenCreated       interface{} `json:"whencreated"`
	Computers         int         `json:"computers"`
	OUs               int         `json:"ous"`
	Users             int         `json:"users"`
	IsACLProtected    bool        `json:"isaclprotected"`
	GPCPath           string      `json:"gpcpath"`
	IsTierZero        bool        `json:"tierzero"`
}

// OU represents a BloodHound OU object.
type OU struct {
	Name              string      `json:"name"`
	ObjectID          string      `json:"objectid"`
	ObjectType        string      `json:"object_type"`
	DistinguishedName string      `json:"distinguishedname"`
	WhenCreated       interface{} `json:"whencreated"`
	Computers         int         `json:"computers"`
	GPOs              int         `json:"gpos"`
	Groups            int         `json:"groups"`
	Users             int         `json:"users"`
	IsACLProtected    bool        `json:"isaclprotected"`
	LastCollected     string      `json:"lastcollected"`
	LastSeen          string      `json:"lastseen"`
	ObjectCategory    string      `json:"objectcategory"`
	WhenChanged       int64       `json:"whenchanged"`
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
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// EntityAdminsResponse wraps a list of entity admins.
type EntityAdminsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// Session represents a user session on a computer.
type Session struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// SessionsResponse wraps a list of sessions.
type SessionsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// Privilege represents a principal with a specific privilege on an entity.
type Privilege struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// PrivilegesResponse wraps a list of privileges.
type PrivilegesResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// ConstrainedDelegation represents a constrained delegation privilege.
type ConstrainedDelegation struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ConstrainedDelegationsResponse wraps a list of constrained delegations.
type ConstrainedDelegationsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// GroupMembership represents a group that a principal is a member of.
type GroupMembership struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// GroupMembershipsResponse wraps a list of group memberships.
type GroupMembershipsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// Controller represents a principal that controls another entity.
type Controller struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ControllersResponse wraps a list of controllers.
type ControllersResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// Controllable represents an entity that is controlled by another principal.
type Controllable struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ControllablesResponse wraps a list of controllables.
type ControllablesResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// ComputersResponse wraps a list of computers.
type ComputersResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
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
	Count int            `json:"count"`
	Limit int            `json:"limit"`
	Skip  int            `json:"skip"`
	Data  []SearchResult `json:"data"`
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

// Container represents a container in BloodHound.
type Container struct {
	Name              string `json:"name"`
	ObjectID          string `json:"objectid"`
	ObjectType        string `json:"object_type"`
	DistinguishedName string `json:"distinguishedname"`
	WhenCreated       interface{} `json:"whencreated"`
}

// DCSyncer represents a principal with DCSync rights.
type DCSyncer struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// OUsResponse wraps a list of OUs.
type OUsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// GroupsResponse wraps a list of groups.
type GroupsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}

// GPOsResponse wraps a list of GPOs.
type GPOsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
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
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
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

// FileUploadJob represents a file upload job.
type FileUploadJob struct {
	ID         int       `json:"id"`
	UserID     string    `json:"user_id"`
	User       User      `json:"user"`
	Status     int       `json:"status"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	LastIngest time.Time `json:"last_ingest"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FileUploadResponse wraps the response from a file upload.
type FileUploadResponse struct {
	Data FileUploadJob `json:"data"`
}

// DomainTrust represents a trust relationship between domains.
type DomainTrust struct {
	Name       string `json:"name"`
	ObjectID   string `json:"objectID"`
	ObjectType string `json:"label"`
	IsTierZero bool   `json:"is_tier_zero"`
}

// ErrorDetail represents a single error in an API response.
type ErrorDetail struct {
	Context string `json:"context"`
	Message string `json:"message"`
}

// ErrorResponse represents the error response from the API.
type ErrorResponse struct {
	HTTPStatus int           `json:"http_status"`
	Timestamp  string        `json:"timestamp"`
	RequestID  string        `json:"request_id"`
	Errors     []ErrorDetail `json:"errors"`
}

// DeleteDatabaseRequest is the payload for deleting the database.
type DeleteDatabaseRequest struct {
	DeleteCollectedGraphData bool `json:"deleteCollectedGraphData"`
	DeleteFileIngestHistory  bool `json:"deleteFileIngestHistory"`
	DeleteDataQualityHistory bool `json:"deleteDataQualityHistory"`
}

// DomainTrustsResponse wraps a list of domain trusts.
type DomainTrustsResponse struct {
	Count int             `json:"count"`
	Limit int             `json:"limit"`
	Skip  int             `json:"skip"`
	Data  json.RawMessage `json:"data"`
}