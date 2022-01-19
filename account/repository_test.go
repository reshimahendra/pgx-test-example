/*
    package account
    repository_test.go 
    - test validity and verify the database operation performed has meet the expectation 
*/
package account

import (
	"errors"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

var (
    // prepare mock
    colums = []string{"id","firstname","lastname","email","passkey"}

    // expected
    want = &User{
        ID : 1,
        Firstname: "John",
        Lastname : "Doe",
        Email : "john@doe.com",
        PassKey: "secret",
    }

)

// Run will prepare our pgxmock connection interface
func Run(t *testing.T) pgxmock.PgxPoolIface{
    t.Helper()
    mock, err := pgxmock.NewPool()
    if err != nil {
        t.Errorf("unexpected error occur: %v\n", err)
    }
    defer mock.Close()

    return mock
}

// TestCreate will test our Create user method
func TestCreate(t *testing.T) {
    mock := Run(t)
    q := `INSERT INTO users (firstname,lastname,email,passkey) 
          VALUES ($1,$2,$3,$4) RETURNING id,firstname,lastname,email,passkey`
    
    // Success
    t.Run("SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(want.Firstname,want.Lastname,want.Email,want.PassKey).
            WillReturnRows(mock.NewRows(colums).
                AddRow(want.ID,want.Firstname,want.Lastname,want.Email,want.PassKey))

        // actual
        ops := NewDatabase(mock)
        got, err := ops.Create(*want)

        assert.NoError(t, err)
        assert.NotNil(t, got)
        assert.Equal(t, got.Firstname, want.Firstname)
        assert.Equal(t, got, want)
    })

    // Test expecting fail/ error
    t.Run("EXPECT FAIL", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WillReturnError(errors.New("error inserting user record"))

        // actual
        ops := NewDatabase(mock)
        got, err := ops.Create(*want)

        // validation & verification
        assert.Error(t, err)
        assert.Nil(t, got)
    })

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectation: %v\n", err)
    }
}

// TestGet will test get one data from database
func TestGet(t *testing.T) {
    mock := Run(t)
    q := `SELECT * FROM users WHERE id = $1`

    t.Run("EXPECT SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(1).
            WillReturnRows(mock.NewRows(colums).
            AddRow(1, "John", "Doe", "john@doe.com", "secret"))

        // actual
        ops := NewDatabase(mock)
        got, err := ops.Get(1) 

        assert.NoError(t, err)
        assert.Equal(t, got.ID, want.ID)
        assert.Equal(t, got, want)
    })

    t.Run("EXPECT FAIL", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(1).
            WillReturnError(errors.New("error getting user data"))

        // actual
        ops := NewDatabase(mock)
        got, err := ops.Get(1) 

        assert.Error(t, err)
        assert.Nil(t, got)
        // assert.NotEqual(t, got, want)
    })

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectation: %v\n", err)
    }
}

// TestGets will test Gets all user data method
func TestGets(t *testing.T) {
    // prepare mock
    mock := Run(t)

    // query to get all user data
    q := `SELECT * FROM users`
    
    // for success test
    users := []*User{
        {1, "john", "doe", "john@doe.com", "secret"},
        {2, "jhonny", "the snail", "jhonny@snail.com", "cretse"},
        {3, "donny", "trumpy", "donny@trumpy", "nohair"},
    }

    // SUCCESS test
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
        WillReturnRows(mock.NewRows(colums).
            AddRow(
                users[0].ID, users[0].Firstname,users[0].Lastname,
                users[0].Email, users[0].PassKey,
            ).
            AddRow(
                users[1].ID, users[1].Firstname,users[1].Lastname,
                users[1].Email, users[1].PassKey,
            ).
            AddRow(
                users[2].ID, users[2].Firstname,users[2].Lastname,
                users[2].Email, users[2].PassKey,
            ),
        )

        ops := NewDatabase(mock)
        got, err := ops.Gets()

        assert.NoError(t, err)
        assert.NotNil(t, got)
        assert.Equal(t, got[0], users[0])
        assert.Equal(t, got, users)
    })

    // FAIL test
    t.Run("EXPECT FAIL", func(t *testing.T){
        // t.Skip()
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WillReturnError(errors.New("error getting user data"))
       

        ops := NewDatabase(mock)
        got, err := ops.Gets()

        // verify and validate
        assert.NotNil(t, ops)
        assert.Error(t, err)
        assert.Nil(t, got)
    })

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectation: %v\n", err)
    }
}

func TestUpdate(t *testing.T) {
    mock := Run(t)
    q := `UPDATE users SET 
            firstname = $2,
            lastname  = $3,
            email = $4,
            passkey = $5
          WHERE id = $1
          RETURNING id, firstname, lastname, email, passkey;
         `

    // SUCCESS test
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey).
            WillReturnRows(mock.NewRows(colums).
            AddRow(want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey),
        )

        ops := NewDatabase(mock)
        got, err := ops.Update(want.ID, *want)

        // t.Logf("GOT :%v, ERR: %v", got, err)
        assert.NotNil(t, ops)
        assert.NoError(t, err)
        assert.NotNil(t, got)
        assert.Equal(t, got.ID, want.ID)
        assert.Equal(t, got, want)
    })

    // FAIL test
    t.Run("EXPECT FAIL", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey).
            WillReturnError(errors.New("update user error"))

        ops := NewDatabase(mock)
        got, err := ops.Update(want.ID, *want)

        // t.Logf("GOT :%v, ERR: %v", got, err)
        assert.NotNil(t, ops)
        assert.Error(t, err)
        assert.Nil(t, got)
    })

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectation: %v\n", err)
    }
}

// TestDelete will test the Delete method of our repository
func TestDelete(t *testing.T) {
    mock := Run(t)

    // query for deleting user data
    // q := `DELETE FROM users WHERE id = $1 RETURNING id`
    q := `DELETE FROM users WHERE id = $1 RETURNING id,firstname,lastname,email,passkey;`

    // SUCCESS test
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(1).
            WillReturnRows(mock.NewRows(colums).
                AddRow(want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey),
            )

        ops := NewDatabase(mock)
        got, err := ops.Delete(1)

        assert.NoError(t, err)
        assert.NotNil(t, ops)
        assert.NotNil(t, got)
        assert.Equal(t, got, want)

        if err := mock.ExpectationsWereMet(); err != nil {
            t.Errorf("there were unfulfilled expectation: %v\n", err)
        }
    })

    // FAIL Test
    t.Run("EXPECT FAIL", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(1).
            WillReturnError(errors.New("error deleting user"))


        ops := NewDatabase(mock)
        got, err := ops.Delete(3)

        assert.Error(t, err)
        assert.NotNil(t, ops)
        assert.Nil(t, got)
    })
}
