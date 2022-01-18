package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTableName will test the tablename func return
func TestTableName(t *testing.T) {
    u := new(User)
    if u.TableName() != "Users" {
        t.Errorf("expecting 'Users' but got '%s'", u.TableName())
    }
}


func TestModelIsValid(t *testing.T) {
    cases := []struct{
        name    string
        user    User
        wanErr  bool
    }{
        {
            "EXPECT VALID",
            User{1, "jhonny", "botak", "jhonny@botak.com", "rahasia"},
            false,
        },
        {
            "EXPECT INVALID 1",
            User{2, "", "", "", ""},
            true,
        },
        {
            "EXPECT INVALID 2",
            User{2, "jhonny", "botak", "", "rahasia"},
            true,
        },

    }

    // run test and loop through test table
    for _, tt := range cases{
        t.Run(tt.name, func(t *testing.T){
            if tt.wanErr {
                assert.Error(t, tt.user.IsValid())
            } else {
                assert.NoError(t, tt.user.IsValid())
            }
        })
    }
}
