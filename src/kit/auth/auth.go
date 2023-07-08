package auth

import (
	"context"
	"net/http"
	"time"
)

// User represents basic user model
type User struct {
	Id          string     // Id - user ID
	Username    string     // Username - username (login)
	Password    string     // Password - password populated if AuthType implies password
	Type        string     // Type - user type
	FirstName   string     // FirstName - first name
	LastName    string     // LastName - last name
	ActivatedAt *time.Time // ActivatedAt - when activated
	LockedAt    *time.Time // LockedAt - when locked
	Groups      []string   // Groups - assigned groups
	Roles       []string   // Roles - user direct roles
}

// UserRepository interface manages users. Must be implemented on client side
type UserRepository interface {
	// GetByUsername gets user by username
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type TokenConfig struct {
	Secret              string
	ExpirationPeriodSec uint `config:"expiration-period-sec"`
}

type Config struct {
	AccessToken  *TokenConfig `config:"access-token"`
	RefreshToken *TokenConfig `config:"refresh-token"`
}

// Session specifies session object
type Session struct {
	Id             string     // Id - session id
	UserId         string     // UserId - Id of logged user
	Username       string     // Username - username of logged user
	LoginAt        time.Time  // LoginAt - when session logged in
	LogoutAt       *time.Time // LogoutAt - when session logged out
	LastActivityAt time.Time  // LastActivityAt - last session activity
	Roles          []string   // Roles - session roles
}

// SessionToken specifies a session token
type SessionToken struct {
	SessionId             string    // SessionId - session ID
	AccessToken           string    // AccessToken
	AccessTokenExpiresAt  time.Time // AccessTokenExpiresAt - when access token expires
	RefreshToken          string    // RefreshToken
	RefreshTokenExpiresAt time.Time // RefreshToken - when refresh token expires
}

// AuthorizationRequest request to authorize session
type AuthorizationRequest struct {
	SessionId              string                   // SessionId session ID
	AuthorizationResources []*AuthorizationResource // AuthorizationResources requested resources
}

type AuthorizeSession interface {
	// AuthorizeSession checks session permissions
	AuthorizeSession(ctx context.Context, rq *AuthorizationRequest) (bool, error)
	// GetRolesForGroups retrieves roles by groups
	GetRolesForGroups(ctx context.Context, groups []string) ([]string, error)
}

type AuthenticateSession interface {
	// AuthSession verifies session token, returns a session if it's verified
	AuthSession(ctx context.Context, token string) (*Session, error)
}

// LoginRequest - specifies login params when logging with password
type LoginRequest struct {
	Username string // Username - mandatory
	Password string // Password
}

// SessionStorage is a session storage interface that must be implemented on client side
type SessionStorage interface {
	// Get - retrieves session by SID
	Get(ctx context.Context, sid string) (*Session, error)
	// GetByUser - retrieves all user's sessions
	GetByUser(ctx context.Context, uid string) ([]*Session, error)
	// CreateSession creates a new session in store
	CreateSession(ctx context.Context, session *Session) error
	// UpdateLastActivity updates last activity date of the session
	UpdateLastActivity(ctx context.Context, sid string, lastActivity time.Time) error
	// Logout marks session as logged out
	Logout(ctx context.Context, sid string, logoutAt time.Time) error
}

type SessionsService interface {
	// Init initializes service with the config
	Init(cfg *Config)
	// LoginPassword logins user with password
	LoginPassword(ctx context.Context, rq *LoginRequest) (*Session, *SessionToken, error)
	// Logout logs out all user's sessions
	Logout(ctx context.Context, sid string) error
	// AuthSession verifies session token, returns a session if it's verified
	AuthSession(ctx context.Context, token string) (*Session, error)
	// Get retrieves a session by sid
	Get(ctx context.Context, sid string) (*Session, error)
	// GetByUser retrieves sessions by user
	GetByUser(ctx context.Context, userId string) ([]*Session, error)
	// RefreshToken allows to refresh a session token
	RefreshToken(ctx context.Context, refreshToken string) (*SessionToken, error)
}

const (
	// permissions
	AccessR = "r" // R read
	AccessW = "w" // W write
	AccessX = "x" // X execute
	AccessD = "d" // D delete
)

type AuthorizationResource struct {
	// resource name
	Resource string
	// list of requested permissions (R, W, X, D)
	Permissions []string
}

// ResourcePolicy
type ResourcePolicy interface {
	// Resolve determines needed authorization resources for the given request
	Resolve(ctx context.Context, r *http.Request) (*AuthorizationResource, error)
}

// ResourcePolicyManager accumulates mapping between URLs and requested resources and then convert it to Authorization request
type ResourcePolicyManager interface {
	// RegisterResourceMapping maps routeId and resource policies
	RegisterResourceMapping(routeId string, policies ...ResourcePolicy)
	// GetRequestedResources resolves policies and retrieves accumulated resources requested to be authorized
	GetRequestedResources(ctx context.Context, routeId string, r *http.Request) ([]*AuthorizationResource, error)
}
