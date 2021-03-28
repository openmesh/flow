package pg

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	db     *sqlx.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context
	// Datasource name.
	DSN string
	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time
}

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sqlx.Tx
	db  *DB
	now time.Time
}

// NewDB returns a new instance of DB associated with the given datasource name.
func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

// Connect to the database.
func (db *DB) Connect() (err error) {
	db.db, err = sqlx.Connect("postgres", db.DSN)
	if err != nil {
		return err
	}

	if err := db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return err
}

// Close the database connection.
func (db *DB) Close() (err error) {
	return db.db.Close()
}

// beginTx starts a transaction and returns a wrapper Tx type. This type
// provides a reference to the database and a fixed timestamp at the start of
// the transaction. The timestamp allows us to mock time during tests as well.
func (db *DB) beginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

// migrate updates the connected database by running any outstanding migration scripts.
func (db *DB) migrate() error {
	driver, err := postgres.WithInstance(db.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://pg/migrations", "boiler", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			return err
		}
		err = nil
	}
	return nil
}

func getRowByID(ctx context.Context, tx *Tx, dest interface{}, id uuid.UUID, table string) error {
	q := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", table)

	if err := tx.Get(dest, q, id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return flow.Errorf(flow.ENOTFOUND, "'%s' with ID %s could not be found.", table, id)
		}
		return err
	}

	return nil
}

// formatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func formatLimitOffset(limit, page int) string {
	if limit > 0 && page > 1 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, (page-1)*limit)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	}
	return ""
}

func where(clause string, val interface{}) string {
	if val == nil {
		return ""
	}
	return "WHERE " + fmt.Sprintf(clause, val)
}

// hashString applies the SHA256 hashing algorithm to a string
func hashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

