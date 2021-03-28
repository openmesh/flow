package flow

import (
	"context"
	"github.com/google/uuid"
)

// contextKey represents an internal key for adding context fields.
// This is considered best practice as it prevents other packages from
// interfering with our context keys.
type contextKey int

// List of context keys.
// These are used to store request-scoped information.
const (
	// Stores the current logged in user in the context.
	userContextKey = contextKey(iota + 1)
)

// NewContextWithUser returns a new context with the given user ID.
func NewContextWithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userContextKey, id)
}

// UserIDFromContext returns the current logged in user.
func UserIDFromContext(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(userContextKey).(uuid.UUID)
	return id
}

