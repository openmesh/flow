package flow

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// User represents a user in the system. Users are typically created via OAuth
// using the AuthService.
type User struct {
	ID uuid.UUID `json:"-" db:"id"`

	// User's preferred name & email.
	Name  *string `json:"name" db:"name"`
	Email *string `json:"email" db:"email"`

	PasswordHash *string `json:"-" db:"password_hash"`

	// Randomly generated API key for use with the CLI.
	APIKey string `json:"-" db:"api_key"`

	// Timestamps for user creation & last update.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// List of associated OAuth authentication objects.
	// Currently only GitHub is supported so there should only be a maximum of one.
	Auths []*Auth `json:"auths" db:"-"`
}

// UserService represents a service for managing users.
type UserService interface {
	// Retrieves a user by ID along with their associated auth objects.
	// Returns ENOTFOUND if user does not exist.
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)

	// Retrieves a list of users by filter. Also returns total count of matching
	// users which may differ from returned results if filter.Limit is specified.
	GetUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)

	// Creates a new user.
	CreateUser(ctx context.Context, user *User) error

	// Updates a user object. Returns EUNAUTHORIZED if current user is not
	// the user that is being updated. Returns ENOTFOUND if user does not exist.
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)

	// Permanently deletes a user and all owned dials. Returns EUNAUTHORIZED
	// if current user is not the user being deleted. Returns ENOTFOUND if
	// user does not exist.
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID     *uuid.UUID `json:"id"`
	Email  *string    `json:"email"`
	APIKey *string    `json:"api_key"`

	// Restrict to subset of results.
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
