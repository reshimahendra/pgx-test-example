/*
    package account
    service.go
        service layer for account package. the business logic of the application
*/
package account

// Interface to Account service
type AccountService interface {
    Get(id int) (*UserResponse, error)
    Gets() ([]*UserResponse, error)
    Create(user User) (*UserResponse, error)
    Update(id int, user User) (*UserResponse, error)
    Delete(id int) (*UserResponse, error)
}

// accountService is wrapper for Database struct
type accountService struct {
    db Database
}

// NewAccountService will create accountService instance
func NewAccountService(db Database) *accountService{
    return &accountService{db: db}
}

// Create method will send create record request to datastore/ repository
func (s *accountService) Create(user User) (*UserResponse, error) {
    // call Create from repository/ datasstore
    u, err := s.db.Create(user)

    // if error occur, return nil rfor the response as well as return the error
    if err != nil {
        return nil, err
    }

    return UserToUserResponse(*u), nil
}

// Get method will get user record by id from repository/ datastore
func (s *accountService) Get(id int) (*UserResponse, error) {
    // call Get from repository/ datastore
    user, err := s.db.Get(id)

    // if error occur, return nil rfor the response as well as return the error
    if err != nil {
        return nil, err
    }

    // return the user response DTO and nil for the error
    return UserToUserResponse(*user), nil 
}

// Gets method will get all user record from repository/ datastore
func (s *accountService) Gets() ([]*UserResponse, error) {
    // Call Gets from repository/ datastore to retreive all User record
    users, err := s.db.Gets()

    // if error occur, return nil for the response slice as well as return the error
    if err != nil {
        return nil, err
    }

    // if no error found, convert all 'User' record to UserResponse dto 
    var uRes []*UserResponse
    for _, user := range users {
        uRes = append(uRes, UserToUserResponse(*user))
    }

    // return response slice and nil if no error found
    return uRes, nil
}

// Update will send update request to datastore/ repository
func (s *accountService) Update(id int, user User) (*UserResponse, error) {
    // check if user data is valid
    // field 'firstname', 'email', and 'passkey' is required
    if err := user.IsValid(); err != nil {
        return nil, err
    }

    // call Update method from repository/ datastore to update certain record
    u, err := s.db.Update(id, user)

    // return nil and the error if error occur
    if err != nil {
        return nil, err
    }

    // return user response dto and nil for the error
    return UserToUserResponse(*u), nil
}

// Delete method will send request to delete record to datastore/ repository
// based on user 'id'
func (s *accountService) Delete(id int) (*UserResponse, error) {
    // call Delete method from repository/ datastore
    u, err := s.db.Delete(id)

    // check if error occur while executing Delete method
    if err != nil {
        return nil, err
    }

    // return user response dto and nil if no error found
    return UserToUserResponse(*u), nil
}

// UserResponse is to response the client/request with 'user' data
type UserResponse struct {
    ID        int       `json:"id"`
    Firstname string    `json:"first_name"`
    Lastname  string    `json:"last_name,omitempty"`
    Email     string    `json:"email"`
    /*
    // must be hidden and not exposed
    password string 
    */
}

// convert 'User' model to 'UserResponse' DTO
func UserToUserResponse(u User) *UserResponse{
    return &UserResponse{
        ID : u.ID,
        Firstname : u.Firstname,
        Lastname : u.Lastname,
        Email : u.Email,
    }
}
