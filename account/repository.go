/*
   package account
   repository.go
   - database operation (CRUD) for account package. its the heart of the apps
*/
package account

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// PgxIface is pgx interface
type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Close(context.Context) error
}

// Database is wrapper for PgxIface
type Database struct {
    DB PgxIface
}

// NewSelector is an initializer for Selector
func NewDatabase(ds PgxIface) Database {
    return Database{DB: ds}
}

// Create method will insert new record to database. 'C' part of the CRUD
func (pool Database) Create(user User) (*User, error) {
    // sql for inserting new record
    q := `INSERT INTO users (id,firstname,lastname,email,passkey)
          VALUES ($1,$2,$3,$4,$5) RETURNING id,firstname,lastname,email,passkey`

    // execute query to insert new record. it takes 'user' variable as its input
    // the result will be placed in 'row' variable
    row := pool.DB.QueryRow(context.Background(), q, 
        user.ID, user.Firstname, user.Lastname, user.Email, user.PassKey)

    // create 'u' variable as 'User' type to contain scanned data value from 'row' variable
    u := new(User)

    // scan 'row' variable and place the value to 'u' variable as well as check for error
    err := row.Scan(
        &u.ID,
        &u.Firstname,
        &u.Lastname,
        &u.Email,
        &u.PassKey,
    )

    // return nil and error if scan operation is fail/ error found
    if err != nil {
        return nil, err
    }

    // return 'u' and nil if no error found
    return u, nil
}

// Get method will get user data by its ID. 'R' part of the CRUD
func (pool Database) Get(id int) (*User, error) {
    // sql command to get user record based on its id
    q := `SELECT * FROM users WHERE id = $1`

    // execute query and place it return value on 'row' variable
    row := pool.DB.QueryRow(context.Background(), q, id)

    // create 'u' variable as User type which will be used as container for
    // 'row' values 
    u := new(User)

    // scan row values and place it in 'u' variable
    err := row.Scan(
        &u.ID,
        &u.Firstname,
        &u.Lastname,
        &u.Email,
        &u.PassKey,
    )

    // return nil and error if error occur while performing 'scan' operation
    if err != nil {
        return nil, err 
    }

    // return 'u' variable and nil if no error found while executing the method
    return u, nil
}

// Gets method will get all user data. extended 'R' part of the CRUD
func (pool Database) Gets() ([]*User, error) {
    // sql comand for getting all user data
    q := `SELECT * FROM users`

    // execute query
    rows, err := pool.DB.Query(context.Background(), q)

    // check if any error occur while executing the query
    if err != nil {
        return nil, err
    }

    // close rows if error ocur
    defer rows.Close()

    // iterate Rows
    var users []*User
    if rows != nil {
        for rows.Next() {
            // create 'u' for struct 'User'
            u := new(User)

            // scan rows and place it in 'u' (user) container
            err := rows.Scan(
                &u.ID,
                &u.Firstname,
                &u.Lastname,
                &u.Email,
                &u.PassKey,
            )

            // return nil and error if scan operation fail
            if err!= nil {
                return nil, err
            }

            // add u to users slice
            users = append(users, u)
        }
    }

    // return users slice and nil for the error
    return users, nil
}

// Update will update user record based on their id
func (pool Database) Update(id int, user User) (*User, error) {
    // prepare update query
    q := `UPDATE users SET 
            firstname = $2,
            lastname  = $3,
            email = $4,
            passkey = $5
          WHERE id = $1
          RETURNING id, firstname, lastname, email, passkey;
         `
    // execute update query
    row := pool.DB.QueryRow(context.Background(), q, id, 
        user.Firstname, user.Lastname, user.Email, user.PassKey)
    
    // create container variable for User
    u := new(User)

    // scan data and place it on 'u' variable we create before and check for error
    if err := row.Scan(
        &u.ID,
        &u.Firstname,
        &u.Lastname,
        &u.Email,
        &u.PassKey,
    ); err != nil {
        return nil, err
    }

    // return variable 'u' as User and nil/ no error
    return u, nil
}

// Delete method will delete user record based on its 'id'
func (pool Database) Delete(id int) (*User, error) {
    // query for deleting user data
    q := `DELETE FROM users WHERE id = $1 RETURNING id,firstname,lastname,email,passkey;`
    
    // execute query
    row := pool.DB.QueryRow(context.Background(), q, id)

    // create container variable for User
    u := new(User)

    // if error occur, return the error
    if err := row.Scan(
        &u.ID,
        &u.Firstname,
        &u.Lastname,
        &u.Email,
        &u.PassKey,
    ); err != nil {
        return nil, err
    }

    // return nil if no error found
    return u, nil

}
