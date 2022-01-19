/*
    package account
    model.go 
        database table representation object model/ object wrapper
*/
package account

import (
	"errors"
)

type IUser interface {
    Get(id int) (*User, error)
    Gets() (Users, error)
    Create(user User) error
    Update(id int, user User) error
}

type Users []*User

// User is user model reflect the 'users' database table
type User struct {
    ID int
    Firstname string
    Lastname string
    Email string
    PassKey string
}

// TableName method will return constant string "Users" as its result
func (u *User) TableName() string {
    return "Users"
}

// IsValid is to validate user input
func (u *User) IsValid() error {
    if u.ID == 0 ||
        u.Firstname == "" ||
        u.Email == "" ||
        u.PassKey == "" {
            return errors.New("user data invalid")
        }

    return nil
}
