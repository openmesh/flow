package pg

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow"
	"io"
)

func getUserByEmail(ctx context.Context, tx *Tx, email string) (*flow.User, error) {
	u, _, err := getUsers(ctx, tx, flow.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	}
	if len(u) == 0 {
		return nil, &flow.Error{Code: flow.ENOTFOUND, Message: "User not found."}
	}
	return u[0], nil
}

func getUsers(ctx context.Context, tx *Tx, filter flow.UserFilter) ([]*flow.User, int, error) {
	var where []string
	var args []interface{}

	if v := filter.ID; v != nil {
		where = append(where, fmt.Sprintf("id = $%d", len(where)+1))
		args = append(args, *v)
	}
	if v := filter.Email; v != nil {
		where = append(where, fmt.Sprintf("email = $%d", len(where)+1))
		args = append(args, *v)
	}
	if v := filter.Email; v != nil {
		where = append(where, fmt.Sprintf("api_key = $%d", len(where)+1))
		args = append(args, *v)
	}

	baseQuery := fmt.Sprintf("SELECT * FROM users %s", buildWhereClause(where))

	var n int
	err := tx.Get(
		&n,
		fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count;", baseQuery),
	)
	if err != nil {
		return nil, n, err
	}

	query := baseQuery + `
		ORDER BY created_at ASC
	` + formatLimitOffset(filter.Limit, filter.Offset)

	users := make([]*flow.User, 0)
	if err := tx.Select(&users, query); err != nil {
		return users, n, err
	}

	return users, n, nil
}

func createUser(ctx context.Context, tx *Tx, user *flow.User) error {
	// Perform basic field validation.
	//if err := user.Validate(); err != nil {
	//	return err
	//}

	// Generate random API key.
	apiKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, apiKey); err != nil {
		return err
	}
	user.APIKey = hex.EncodeToString(apiKey)

	// Execute insertion query.
	stmt, err := tx.PrepareNamed(`
		INSERT INTO
			users 
				(
				 	name,
				 	email,
				 	api_key
				)
		VALUES 
			(
			 	:name,
			 	:email,
			 	:api_key
			)
		RETURNING 
			*
	`)
	if err != nil {
		return err
	}
	var res flow.User
	if err := stmt.Get(&res, user); err != nil {
		return err
	}
	user = &res

	return nil
}

// getUserByID is a helper function to fetch a user by ID.
// Returns ENOTFOUND if user does not exist.
func getUserByID(ctx context.Context, tx *Tx, id uuid.UUID) (*flow.User, error) {
	a, _, err := getUsers(ctx, tx, flow.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &flow.Error{Code: flow.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}