/*
    package account
    db.go
    - holding the main database configuration and setup for the package
*/
package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DatabaseConfig is configuration wrapper for our database setting
type DatabaseConfig struct {
    Username string
    Password string
    Hostname string
    Port string
    DBName string
}

// DSN will get datasource name of the database configuration
func (db DatabaseConfig) DSN() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", 
        db.Username, db.Password, db.Hostname, db.Port, db.DBName)
}

// Datastore is pgxpool.Pool wrapper
type Datastore struct{
    dbPool *pgxpool.Pool
}

// NewDatastore will create datastore instance
func NewDatastore(pool *pgxpool.Pool) Datastore {
    return Datastore{dbPool: pool}
}

// Pool method will get current pgxpool.Pool pointer
func (ds Datastore) Pool() *pgxpool.Pool{
    return ds.dbPool
}

// NewDBPool will create pool connection to database
func NewDBPool(dbConfig DatabaseConfig) (*pgxpool.Pool, func(), error) {

	f := func() {}

    // create pgx connection pool
	pool, err := pgxpool.Connect(context.Background(), dbConfig.DSN())

    // return nil to connection and return error if error occur
	if err != nil {
        return nil, f, errors.New("database connection error")
	}


    // validateDBPool
	err = validateDBPool(pool)

    // return nil and error if error occur
	if err != nil {
		return nil, f, err
	}

    // return connection pool and inline function to close/ clear the pool if not used. 
    // return nil for the error since there should be no error to this point
	return pool, func() { pool.Close() }, nil
}

// validateDBPool will pings the database and logs the current user and database
func validateDBPool(pool *pgxpool.Pool) error {
	// tried to ping connection
    err := pool.Ping(context.Background())

    // return error if error found
	if err != nil {
        return errors.New("database connection error")
	}

	var (
		currentDatabase string
		currentUser     string
		dbVersion       string
	)
	
    // Lets try to get db system info
    sqlStatement := `select current_database(), current_user, version();`
	row := pool.QueryRow(context.Background(), sqlStatement)
	err = row.Scan(&currentDatabase, &currentUser, &dbVersion)

	switch {
	case err == sql.ErrNoRows:
		return errors.New("no rows were returned")
	case err != nil:
		return errors.New("database connection error")
	default:
		log.Printf("database version: %s\n", dbVersion)
		log.Printf("current database user: %s\n", currentUser)
		log.Printf("current database: %s\n", currentDatabase)
	}

	return nil
}
