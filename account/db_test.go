package account

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

// TestDSN will test DSN method
func TestDSN(t *testing.T) {
    dbconf := &DatabaseConfig{
        Username:"user",
        Password:"pass",
        Hostname:"localhost",
        Port: "5432",
        DBName: "testdb",
    }
    
    got := dbconf.DSN()
    want := "postgres://user:pass@localhost:5432/testdb"
    assert.Equal(t, got, want)
}

// TestNewDatastore is for testing connection pool to database
func TestNewDatastore(t *testing.T) {
    // mock, err := pgxmock.NewPool()
    dbconf := &DatabaseConfig{
        Username:"golang",
        Password:"golang",
        Hostname:"localhost",
        Port: "5432",
        DBName: "golangtest",
    }
    pool, err := pgxpool.Connect(context.Background(), dbconf.DSN())
    if err != nil {
        t.Errorf("error creating database connection pool stub")
    }
    pool.Config().MaxConnIdleTime = time.Duration(time.Second * 1)
    
    defer pool.Close()

    tds := NewDatastore(pool)

    assert.NotNil(t, tds)
    assert.Equal(t, tds.Pool(), pool)
} 

func TestNewDBPool(t *testing.T) {
    db := &DatabaseConfig{
        Username : "golang",
        Password : "golang",
        Hostname : "localhost",
        Port     : "5432",
        DBName   : "golangtest",
    }

    pool, cleanup, err := NewDBPool(*db)
    if err != nil {
        t.Errorf("error creating database connection pool stub")
    }

    defer pool.Close()
    pool.Config().MaxConnIdleTime = time.Duration(time.Second * 1)
    pool.Config().MaxConnLifetime = time.Duration(time.Second * 2)

    t.Cleanup(cleanup)
}
