package flow

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// Authentication providers.
const (
	AuthSourceInternal = "internal"
	AuthSourceGitHub   = "github"
)

// Auth represents a set of credentials. These are linked to a User so a
// single user could authenticate through multiple providers.
//
// The authentication system links users by email address.
type Auth struct {
	ID uuid.UUID `json:"-" db:"id"`

	// User can have one or more methods of authentication.
	// However, only one per source is allowed per user.
	UserID uuid.UUID `json:"-" db:"user_id"`
	User   *User     `json:"user,omitempty" db:"-"`

	// The authentication source & the source provider's user ID.
	// Source can only be "github" currently.
	Source   string `json:"source" db:"source"`
	SourceID string `json:"-" db:"source_id"`

	// OAuth fields returned from the authentication provider.
	// GitHub does not use refresh tokens but the field exists for future providers.
	AccessToken  string     `json:"-" db:"access_token"`
	RefreshToken string     `json:"-" db:"refresh_token"`
	ExpiresAt    *time.Time `json:"-" db:"expires_at"`

	// Timestamps of creation & last update.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AuthService represents a service for managing auths.
type AuthService interface {
	// Looks up an authentication object by ID along with the associated user.
	// Returns ENOTFOUND if ID does not exist.
	GetAuthByID(ctx context.Context, id uuid.UUID) (*Auth, error)

	// Retrieves authentication objects based on a filter. Also returns the
	// total number of objects that match the filter. This may differ from the
	// returned object count if the Limit field is set.
	GetAuths(ctx context.Context, filter AuthFilter) ([]*Auth, int, error)

	// Creates a new authentication object. If a User is attached to auth, then
	// the auth object is linked to an existing user. Otherwise a new user
	// object is created.
	//
	// On success, the auth.ID is set to the new authentication ID.
	CreateAuth(ctx context.Context, auth *Auth) error

	// Permanently deletes an authentication object from the system by ID.
	// The parent user object is not removed.
	DeleteAuth(ctx context.Context, id uuid.UUID) error

	SignUp(ctx context.Context, email string, name string, password string) (*User, error)
}

// AuthFilter represents a filter accepted by FindAuths().
type AuthFilter struct {
	// Filtering fields.
	ID       *uuid.UUID `json:"id"`
	UserID   *int       `json:"user_id"`
	Source   *string    `json:"source"`
	SourceID *string    `json:"source_id"`

	// Restricts results to a subset of the total range.
	// Can be used for pagination.
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
