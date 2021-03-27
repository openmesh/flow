package pg

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow"
	"strings"
	"time"
)

type authService struct {
	db *DB
}

func NewAuthService(db *DB) flow.AuthService {
	return authService{db}
}

func (s authService) GetAuthByID(ctx context.Context, id uuid.UUID) (*flow.Auth, error) {
	tx, err := s.db.beginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	auth, err := getAuthByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	if err := attachAuthAssociations(ctx, tx, auth); err != nil {
		return nil, err
	}

	return auth, nil
}

func (s authService) GetAuths(ctx context.Context, filter flow.AuthFilter) ([]*flow.Auth, int, error) {
	tx, err := s.db.beginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}

	// Fetch the individual auth entities from the database.
	auths, n, err := getAuths(ctx, tx, filter)
	if err != nil {
		return auths, n, err
	}

	// TODO figure out how to batch these queries
	for _, auth := range auths {
		if err := attachAuthAssociations(ctx, tx, auth); err != nil {
			return auths, n, err
		}
	}

	return auths, n, nil
}

func (s authService) CreateAuth(ctx context.Context, auth *flow.Auth) error {
	tx, err := s.db.beginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check to see if the auth already exists for the given source.
	if other, err := getAuthBySourceID(ctx, tx, auth.Source, auth.SourceID); err == nil {
		// If an auth already exists for the source user, update with the new tokens.
		if other, err = updateAuth(ctx, tx, other.ID, auth.AccessToken, auth.RefreshToken, auth.ExpiresAt); err != nil {
			return fmt.Errorf("cannot update auth: id=%d err=%w", other.ID, err)
		} else if err := attachAuthAssociations(ctx, tx, other); err != nil {
			return err
		}

		// Copy found auth back to the caller's arg & return.
		*auth = *other
		return tx.Commit()
	} else if flow.ErrorCode(err) != flow.ENOTFOUND {
		return fmt.Errorf("canot find auth by source user: %w", err)
	}

	// Check if auth has a new user object passed in. It is considered "new" if
	// the caller doesn't know the database ID for the user.
	if auth.UserID == uuid.New() && auth.User != nil {
		// Look up the user by email address. If no user can be found then
		// create a new user with the auth.User object passed in.
		if user, err := getUserByEmail(ctx, tx, auth.User.Email); err == nil { // user exists
			auth.User = user
		} else if flow.ErrorCode(err) == flow.ENOTFOUND { // user does not exist
			if err := createUser(ctx, tx, auth.User); err != nil {
				return fmt.Errorf("cannot create user: %w", err)
			}
		} else {
			return fmt.Errorf("cannot find user by email: %w", err)
		}

		// Assign the created/found user ID back to the auth object.
		auth.UserID = auth.User.ID
	}

	// Create new auth object & attach associated user.
	if err := createAuth(ctx, tx, auth); err != nil {
		return err
	} else if err := attachAuthAssociations(ctx, tx, auth); err != nil {
		return err
	}
	return tx.Commit()
}

func (s authService) DeleteAuth(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

// getAuthByID is a helper function that returns an auth object by ID.
// Returns ENOTFOUND if auth doesn't exist.
func getAuthByID(ctx context.Context, tx *Tx, id uuid.UUID) (*flow.Auth, error) {
	auths, _, err := getAuths(ctx, tx, flow.AuthFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	if len(auths) == 0 {
		return nil, &flow.Error{
			Code:    flow.ENOTFOUND,
			Message: "Auth not found.",
		}
	}
	return auths[0], nil
}

// getAuthBySourceID is a helper function to return an auth object by source ID.
// Returns ENOTFOUND if auth doesn't exist.
func getAuthBySourceID(ctx context.Context, tx *Tx, source, sourceID string) (*flow.Auth, error) {
	auths, _, err := getAuths(ctx, tx, flow.AuthFilter{Source: &source, SourceID: &sourceID})
	if err != nil {
		return nil, err
	} else if len(auths) == 0 {
		return nil, &flow.Error{Code: flow.ENOTFOUND, Message: "Auth not found."}
	}
	return auths[0], nil
}

// getAuths returns a list of auth objects that match a filter. Aso returns a total count of matches
// which may differ results length if filter.Limit is set.
func getAuths(ctx context.Context, tx *Tx, filter flow.AuthFilter) ([]*flow.Auth, int, error) {
	var where []string
	var args []interface{}

	if v := filter.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", len(where)+1))
		args = append(args, *v)
	}
	if v := filter.UserID; v != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", len(where)+1))
		args = append(args, *v)
	}
	if v := filter.Source; v != nil {
		where = append(where, fmt.Sprintf("source = $%d", len(where)+1))
		args = append(args, *v)
	}
	if v := filter.SourceID; v != nil {
		where = append(where, fmt.Sprintf("source_id = $%d", len(where)+1))
		args = append(args, *v)
	}

	baseQuery := fmt.Sprintf("SELECT * FROM auths %s", buildWhereClause(where))

	var n int
	err := tx.Get(
		&n,
		fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count;", baseQuery),
		args...,
	)
	if err != nil {
		return nil, n, err
	}

	query := baseQuery + `
		ORDER BY created_at ASC
		` + formatLimitOffset(filter.Limit, filter.Offset)

	auths := make([]*flow.Auth, 0)
	if err := tx.Select(&auths, query); err != nil {
		return auths, n, err
	}

	return auths, n, nil
}

func createAuth(ctx context.Context, tx *Tx, auth *flow.Auth) error {
	stmt, err := tx.PrepareNamed(`
		INSERT INTO 
			auths
				(
					user_id, 
				 	source,
				 	source_id,
				 	access_token,
				 	refresh_token,
				 	expires_at
				)
		VALUES 
			(
				:user_id, 
				:source,
				:source_id,
				:access_token,
				:refresh_token,
				:expires_at
			)
	`)
	if err != nil {
		return err
	}

	var res flow.Auth
	if err := stmt.Get(&res, auth); err != nil {
		return err
	}
	auth = &res

	return nil
}

func updateAuth(ctx context.Context, tx *Tx, id uuid.UUID, accessToken, refreshToken string, expiresAt *time.Time) (*flow.Auth, error) {
	// Fetch currency entity state
	auth, err := getAuthByID(ctx, tx, id)
	if err != nil {
		return auth, err
	}

	// Update fields
	auth.AccessToken = accessToken
	auth.RefreshToken = refreshToken
	auth.ExpiresAt = expiresAt

	// TODO add validation
	//if err := auth.Validate(); err != nil {
	//	return auth, err
	//}

	if _, err := tx.ExecContext(ctx, `
		UPDATE
		    auths
		SET
			access_token = $1,
		    refresh_token = $2,
		    expires_at = $3,
		WHERE id = $4
	`,
		args...,
	); err != nil {
		return nil, err
	}

	return auth, nil
}

func buildWhereClause(clauses []string) string {
	if len(clauses) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(clauses, " AND ")
}

// attachAuthAssociations is a helper function to fetch and attach the associated user to an auth
// object.
func attachAuthAssociations(ctx context.Context, tx *Tx, auth *flow.Auth) error {
	var err error
	if auth.User, err = getUserByID(ctx, tx, auth.UserID); err != nil {
		return fmt.Errorf("failed to attach auth user: %w", err)
	}
	return nil
}
