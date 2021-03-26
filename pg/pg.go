package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmesh/flow"
	"reflect"
	"strings"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
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

// insert an entity into the database. The reflect package is used to build the query from the provided entity. `entity`
// should be a pointer to the struct that should be inserted into the database. It is expected that the underlying
// struct would have the following properties: ID, CompanyID, CreatedAt, UpdatedAt, CreatedBy, UpdatedBy. These
// fields are set by the method using the information about the user that is extracted from ctx.
func insert(ctx context.Context, tx *Tx, entity interface{}, table string) error {
	// Get current user from context.
	// user := session.UserFromContext(ctx)

	// currTime := time.Now()

	// Set metadata for struct to be inserted into the database.
	// reflect.Indirect(reflect.ValueOf(entity)).FieldByName("CompanyID").SetInt(int64(user.CompanyID))
	// reflect.Indirect(reflect.ValueOf(entity)).FieldByName("CreatedAt").Set(reflect.ValueOf(currTime))
	// reflect.Indirect(reflect.ValueOf(entity)).FieldByName("UpdatedAt").Set(reflect.ValueOf(currTime))
	// reflect.Indirect(reflect.ValueOf(entity)).FieldByName("CreatedBy").SetInt(int64(user.ID))
	// reflect.Indirect(reflect.ValueOf(entity)).FieldByName("UpdatedBy").SetInt(int64(user.ID))

	// Build a slice containing the column names for all fields to be inserted into the database.
	var columns []string

	t := reflect.Indirect(reflect.ValueOf(entity)).Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// If a field has no `db` tag then the struct is invalid and we should fail
		column, ok := field.Tag.Lookup("db")
		if !ok {
			return flow.Errorf(flow.EINTERNAL,
				"Field '%s' does not contain a `db` tag.", field.Name)
		}
		// The tag `id,omitempty` should be used for primary key fields and should be omitted from the query. The tag `-`
		// should be used for related structs and properties that do not have an underlying value persisted in the database.
		// These values should also be omitted from the query.
		if column == "id,omitempty" || column == "-" {
			continue
		}
		columns = append(columns, column)
	}

	// Build the SQL query to be executed from the column names that have been derived from the struct's tags and the
	// `table` argument.
	q := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id;",
		table,
		strings.Join(columns, ", "),
		":"+strings.Join(columns, ", :"),
	)

	// Use the sqlx package to execute the query.
	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return err
	}

	var id uuid.UUID
	err = stmt.Get(&id, entity)
	if err != nil {
		return err
	}

	// Assign the returned ID value to the entity.
	reflect.Indirect(reflect.ValueOf(entity)).FieldByName("ID").Set(reflect.ValueOf(id))

	return nil
}

