/*
    package account
    service_test.go 
    - test validity and verify the business layer operation to meet the expectation 
*/
package account

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

// Setup will prepare mock and account service instance
func Setup(t *testing.T) (pgxmock.PgxConnIface, accountService){
    t.Helper()
    mock, err := pgxmock.NewConn()

    if err != nil {
        t.Errorf("error creating stub connection: %v\n", err)
    }
    defer mock.Close(context.Background())

    db := NewDatabase(mock)
    svc := NewAccountService(db)

    assert.NotNil(t, db)
    assert.NotNil(t, svc)

    return mock, *svc
}

// TestAccountServiceCreate to test Create service from account service
func TestAccountServiceCreate(t *testing.T) {
    // prepare mock and service
    mock, service := Setup(t)

    // sql for inserting new record
    q := `INSERT INTO users (id,firstname,lastname,email,passkey)
          VALUES ($1,$2,$3,$4,$5) RETURNING id,firstname,lastname,email,passkey`

    // SUCCESS test
    t.Run("EXPECT SUCCESS", func(t *testing.T) {
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey).
            WillReturnRows(pgxmock.NewRows(colums).
                AddRow(want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey),
            )

        // actual
        got, err := service.Create(*want)

        // validation and verification
        assert.NoError(t, err)
        assert.NotNil(t, got)
        assert.Equal(t, got, UserToUserResponse(*want))
    })

    // EXPECT FAIL test
    t.Run("EXPECT FAIL", func(t *testing.T) {
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(want.ID, "", want.Lastname, want.Email, want.PassKey).
            WillReturnError(errors.New("error inserting user record"))

        // actual
        got, err := service.Create(*want)

        // validation and verification
        assert.Error(t, err)
        assert.Nil(t, got)
    })
}

// TestAccountServiceGet to test Get Service from AccountService
func TestAccountServiceGet(t *testing.T) {
    mock, service := Setup(t)
    q := `SELECT * FROM users WHERE id = $1`

    // SUCCESS test
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(1).
            WillReturnRows(mock.NewRows(colums).AddRow(
                want.ID, want.Firstname, want.Lastname, want.Email, want.PassKey,
            ))

        // actual
        got, err := service.Get(1)

        // validation and verification
        assert.NoError(t, err)
        assert.Equal(t, got, UserToUserResponse(*want))
    })

    // FAIL test
    t.Run("EXPECT FAIL", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WithArgs(3).
            WillReturnError(errors.New("user not found"))

        // actual
        got, err := service.Get(3)

        // validation and verification
        assert.Error(t, err)
        assert.Nil(t, got)
    })

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectation: %v\n", err)
    }
}

// TestAccountServiceGets to test Gets Service from AccountService
func TestAccountServiceGets(t *testing.T) {
    mock, service := Setup(t)
    q := `SELECT * FROM users`

    // SUCCESS test
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WillReturnRows(mock.NewRows(colums).
                AddRow(1, "john", "doe", "john@doe.com", "secret").
                AddRow(2, "donny", "trumpy", "donny@trumpy.com", "nohair"),
            )

            want := []*UserResponse{
                {ID:1,Firstname:"john",Lastname:"doe",Email:"john@doe.com"},
                {ID:2,Firstname:"donny",Lastname:"trumpy",Email:"donny@trumpy.com"},
            }

        // actual
        got, err := service.Gets()

        // validation and verification
        assert.NoError(t, err)
        assert.Equal(t, got, want)
    })

    // FAIL test
    t.Run("EXPECT FAIL", func(t *testing.T){
        mock.ExpectQuery(regexp.QuoteMeta(q)).
            WillReturnError(errors.New("user not found"))

        // actual
        got, err := service.Gets()

        // validation and verification
        assert.Error(t, err)
        assert.Nil(t, got)
    })

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectation: %v\n", err)
    }
}

// TestUserToUserResponse is to test UserToUserResponse function
func TestUserToUserResponse(t *testing.T) {
    u := User{
        ID : 1,
        Firstname : "john",
        Lastname : "doe",
        Email : "john@doe.com",
        PassKey : "secret", 
    }

    got := UserToUserResponse(u)
    assert.NotNil(t, got)
    assert.Equal(t, got.ID, 1)
    assert.Equal(t, got.Firstname, "john")
    assert.Equal(t, got.Lastname, "doe")
    assert.Equal(t, got.Email, "john@doe.com")
}
