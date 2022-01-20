/*
   package account
   handler_test.go
       test handler/ controller layer for account package.
*/
package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAccService struct {
    t *testing.T
}

var (
    users = []*User{
        {ID: 1, Firstname:"joe",Lastname:"taslim",Email:"joe@taslim.com",PassKey:"secret"},
        {ID: 2, Firstname:"john",Lastname:"doe",Email:"john@doe.com",PassKey:"secret"},
        {ID: 3, Firstname:"janne",Lastname:"doe",Email:"janne@doe.com",PassKey:"secret"},
    }
    wantError bool
)

// usersResponse will convert user slice into userResponse slice
func usersResponse() []*UserResponse {
    var ur []*UserResponse
    for _, u := range users{
        ur = append(ur, UserToUserResponse(*u))
    }
    return ur
}

func NewMockAccService(t *testing.T) *mockAccService{
    return &mockAccService{t:t}
}

// Create method is 'mock' to satisfy 'Create' method for AccountService interface
// its act as the 'double' or as a 'counterfeiter' for AccountService.Create
func (m *mockAccService) Create(user User) (*UserResponse, error) {
    if user.ID == 0 || user.Firstname == "" || user.Email=="" || user.PassKey=="" {
        return nil, errors.New("user invalid")
    }
    return UserToUserResponse(user), nil
}

// Get method is 'mock' to satisfy 'Get' method for AccountService interface
// its act as the 'double' or as a 'counterfeiter' for AccountService.Get
func (m *mockAccService) Get(id int) (*UserResponse, error) {
    if len(users) < id {
        return nil, errors.New("data not found")
    }
    return UserToUserResponse(*users[id-1]), nil
}

// Gets method is 'mock' to satisfy 'Gets' method for AccountService interface
// its act as the 'double' or as a 'counterfeiter' for AccountService.Gets
func (m *mockAccService) Gets() ([]*UserResponse, error) {
    // to generate/ force error return on test
    if wantError {
        return nil, errors.New("error found")
    }

    // normal succes return
    return usersResponse(), nil 
}

// Update method is 'mock' to satisfy 'Update' method for AccountService interface
// its act as the 'double' or as a 'counterfeiter' for AccountService.Update
func (m *mockAccService) Update(id int, user User) (*UserResponse, error) {
    if len(users) < id {
        return nil, errors.New("data not found")
    }
    if user.ID == 0 || user.Firstname == "" || user.Email=="" || user.PassKey=="" {
        return nil, errors.New("user invalid")
    }

    return UserToUserResponse(user), nil
}

// Delete method is 'mock' to satisfy 'Delete' method for AccountService interface
// its act as the 'double' or as a 'counterfeiter' for AccountService.Delete
func (m *mockAccService) Delete(id int) (*UserResponse, error) {
    if len(users) < id {
        return nil, errors.New("data not found")
    }
    return UserToUserResponse(*users[id-1]), nil
}

// NewTestHandler is to prepare needed instance before the test executed
func NewTestHandler(t *testing.T) (accountHandler) {
    t.Helper()

    gin.SetMode(gin.TestMode)

    mock := NewMockAccService(t)
    handler := NewAccountHandler(mock)

    return handler
}

func NewTestRecordWriter() (*httptest.ResponseRecorder, *gin.Context) {
    writer := httptest.NewRecorder()
    context, _ := gin.CreateTestContext(writer)

    return writer, context
}

// TestUserCreateHandler is routine for testing UserCreateHandler
func TestUserCreateHandler (t *testing.T) {
    // prepare the test
    handler := NewTestHandler(t)

    // EXPECT SUCCESS
    // should return 200/ status ok
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        // prepare request/ response / gin context
        writer, context := NewTestRecordWriter()

        // marshal json for 'User' model from 'users' slice index 1
        userJSON, err := json.Marshal(&users[0])
        assert.NoError(t, err)

        // inject json to the request body
        context.Request, err = http.NewRequest("POST", "/", bytes.NewBuffer(userJSON))
        assert.NoError(t, err)

        // make sure to add aplication/json as it content type
        context.Request.Header.Add("content-type", "application/json")

        // actual method handler executed
        handler.UserCreateHandler(context)
        
        // prepare expected data so we can compare with the actual/ result/ response data
        want, err := json.Marshal(UserToUserResponse(*users[0]))
        assert.NoError(t, err)

        // make sure response status and body is equal to the expectation
        assert.Equal(t, http.StatusOK, writer.Code)
        assert.Equal(t, want, writer.Body.Bytes())
        // t.Logf("handler: %v\nwriter:%v\ncontext:%v\n", handler, writer, context)
    })

    // EXPECT FAIL error bind json
    // should return 400/ bad request
    t.Run("EXPECT ERROR bind json fail", func(t *testing.T){
        // prepare request/ response/ context
        writer, context := NewTestRecordWriter()

        // dont include user data to request body (send nil), so binding json will fail
        context.Request, _ = http.NewRequest("POST", "/", nil)
        context.Request.Header.Add("content-type", "application/json")

        // actual method to test
        handler.UserCreateHandler(context)
        
        // expected http status response is 400
        assert.Equal(t, http.StatusBadRequest, writer.Code)
    })

    // EXPECT FAIL error create record
    // should return 500/ internal server error
    t.Run("EXPECT FAIL error create record", func(t *testing.T){
        // prepare request/ response/ context
        writer, context := NewTestRecordWriter()

        // we just input the user ID and ignore the rest including ignoring required field 
        userJSON, err := json.Marshal(User{ID:1})
        assert.NoError(t, err)

        // insert the uncomplete data to request body
        context.Request, err = http.NewRequest("POST", "/", bytes.NewBuffer(userJSON))
        assert.NoError(t, err)
        context.Request.Header.Add("content-type", "application/json")

        // actual method we want to test
        handler.UserCreateHandler(context)
        
        // make sure the http status response as we expected (500)
        assert.Equal(t, http.StatusInternalServerError, writer.Code)
    })
}

// TestUserGetHandler is for testing UserGetHandler behaviour
func TestUserGetHandler(t *testing.T) {
    // prepare test
    handler := NewTestHandler(t)

    // EXPECT SUCCESS
    // should return 200/ status ok
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        // prepare http response/ writer and gin context
        writer, context := NewTestRecordWriter()
        context.Params = gin.Params{
            {Key:"id", Value:"1"},
        }

        // actual method call
        handler.UserGetHandler(context)

        // make sure expected status code equal to response status code
        assert.Equal(t, http.StatusOK, writer.Code)

        // marshal expected value using expected data 
        // to compare with the actual/ response body
        want, err := json.Marshal(UserToUserResponse(*users[0]))
        assert.NoError(t, err)

        // make sure expected body value match with the actual/ response body value
        assert.Equal(t, want, writer.Body.Bytes())
    })

    // EXPECT FAIL param error
    // Should return 400/ bad request status
    // since it expecting fail, assert the body is not required
    t.Run("EXPECT FAIL param empty", func(t *testing.T){
        // prepare http response/ writer and gin context
        writer, context := NewTestRecordWriter()
        
        // actual method to test
        handler.UserGetHandler(context)

        // make sure expected body value match with the actual/ response body value
        assert.Equal(t, http.StatusBadRequest, writer.Code)
    })

    // EXPECT FAIL error get data (out of range index on users)
    // should return 500/ internal server error 
    // since it expecting fail, no need to assert the body
    t.Run("EXPECT FAIL service get error", func(t *testing.T){
        writer, context := NewTestRecordWriter()
        context.Params = gin.Params{
            {Key:"id", Value:"5"},
        }
        handler.UserGetHandler(context)

        // make sure expected body value match with the actual/ response body value
        assert.Equal(t, http.StatusInternalServerError, writer.Code)
    })
} 

// TestUserGetsHandler is for testing UserGetsHandler behaviour
func TestUserGetsHandler(t *testing.T) {
    // prepare test
    handler := NewTestHandler(t)

    // EXPECT SUCCESS
    // should return 200/ status ok
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        // prepare http response/ writer and gin context
        writer, context := NewTestRecordWriter()

        // actual method call
        handler.UserGetsHandler(context)

        // make sure expected status code equal to response status code
        assert.Equal(t, http.StatusOK, writer.Code)

        // marshal expected value using expected data 
        // to compare with the actual/ response body
        want, err := json.Marshal(usersResponse())
        assert.NoError(t, err)

        // make sure expected body value match with the actual/ response body value
        assert.Equal(t, want, writer.Body.Bytes())
    })

    // EXPECT FAIL error get data (empty users)
    // should return 500/ internal server error 
    // since it expecting fail, no need to assert the body
    t.Run("EXPECT FAIL service get error", func(t *testing.T){
        writer, context := NewTestRecordWriter()
        
        // force to return error
        wantError = true

        // actual method to get all user data
        handler.UserGetsHandler(context)

        // make sure expected body value match with the actual/ response body value
        assert.Equal(t, http.StatusInternalServerError, writer.Code)

        // return wantError value to its default in case any routine also
        // use it
        wantError = false
    })
}

// TestUserUpdateHandler will test UserUpdateHandler behaviour
func TestUserUpdateHandler(t *testing.T) {
    // prepare test 
    handler := NewTestHandler(t)

    // EXPECT SUCCESS, return http status 200
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        // prepare response/ writer/ context
        writer, context := NewTestRecordWriter()

        // set request param
        context.Params = gin.Params{
            {Key:"id", Value:"1"},
        }

        // marshal json with new user data
        userJSON, err := json.Marshal(User{
            ID : 1,
            Firstname : "zhao",
            Lastname : "lucy",
            Email : "zhao@lucy.com",
            PassKey : "lucysecret",
        })
        assert.NoError(t, err)

        // inject json to the request body for update operation (PUT)
        context.Request, err = http.NewRequest("PUT", "/", bytes.NewBuffer(userJSON))
        assert.NoError(t, err)

        // make sure to add aplication/json as it content type
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed
        handler.UserUpdateHandler(context)

        // make sure the expected status code (200/ status ok) match with the
        // response status code
        assert.Equal(t, http.StatusOK, writer.Code)
    })

    // EXPECT FAIL bad param, return http status 400/ bad request
    // we'll simulate it by removing the param
    t.Run("EXPECT FAIL error param id", func(t *testing.T){
        // prepare response/ writer/ context
        writer, context := NewTestRecordWriter()

        // marshal json with new user data
        userJSON, err := json.Marshal(User{
            ID : 1,
            Firstname : "zhao",
            Lastname : "lucy",
            Email : "zhao@lucy.com",
            PassKey : "lucysecret",
        })
        assert.NoError(t, err)

        // inject json to the request body for update operation (PUT)
        context.Request, err = http.NewRequest("PUT", "/", bytes.NewBuffer(userJSON))
        assert.NoError(t, err)

        // make sure to add aplication/json as it content type
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed
        handler.UserUpdateHandler(context)

        // make sure the expected status code (400/ status bad request) match with the
        // response status code
        assert.Equal(t, http.StatusBadRequest, writer.Code)
    })

    // EXPECT FAIL json bind fail, return http status 400/ bad request
    // we'll simulate it by removing the request body
    t.Run("EXPECT FAIL error user input", func(t *testing.T){
        // prepare response/ writer/ context
        writer, context := NewTestRecordWriter()

        // set request param
        context.Params = gin.Params{
            {Key:"id", Value:"1"},
        }

        // no json data injected to the request body to simulate bad request error
        context.Request, _ = http.NewRequest("PUT", "/", nil )

        // make sure to add aplication/json as it content type
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed
        handler.UserUpdateHandler(context)

        // make sure the expected status code (400/ status bad request) match with the
        // response status code
        assert.Equal(t, http.StatusBadRequest, writer.Code)
    })

    // EXPECT FAIL error update data, return http status 500/ internal server error
    // we'll simulate it by removing 1 of required field, so it will still pass 
    // json binding
    t.Run("EXPECT FAIL error update data", func(t *testing.T){
        // prepare response/ writer/ context
        writer, context := NewTestRecordWriter()

        // set request param 
        context.Params = gin.Params{
            {Key:"id", Value:"1"},
        }

        // marshal json with new user data
        // to make it fail, we will remove 1 of required field
        userJSON, err := json.Marshal(User{
            ID : 1,
            // Firstname : "zhao",
            Lastname : "lucy",
            Email : "zhao@lucy.com",
            PassKey : "lucysecret",
        })
        assert.NoError(t, err)

        // inject invalid data to the request body to simulate how
        // we can get the internal server error
        context.Request, _ = http.NewRequest("PUT", "/", bytes.NewBuffer(userJSON))

        // make sure to add aplication/json as it content type
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed
        handler.UserUpdateHandler(context)

        // make sure the expected status code (500/ status interna; server error) 
        // match with the response status code
        assert.Equal(t, http.StatusInternalServerError, writer.Code)
    })
}

// TestUserDeleteHandler will simulate and test the behaviour of
// UserDeleteHandler method
func TestUserDeleteHandler(t *testing.T) {
    // prepare the mocked test handler
    handler := NewTestHandler(t)

    // EXPECT SUCCESS will return 200/ status ok
    t.Run("EXPECT SUCCESS", func(t *testing.T){
        // prepare request/ writer/ context 
        writer, context := NewTestRecordWriter()

        // prepare request param id. we will simulate to delete 
        // user id with id = 1
        context.Params = gin.Params{
            {Key: "id", Value:"1"},
        }

        // create new request with 'DELETE' method
        context.Request, _ = http.NewRequest("DELETE", "/", nil)
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed to test it behaviour
        handler.UserDeleteHandler(context)

        // make sure the expected status code (200/ status ok) are match with
        // the response code
        assert.Equal(t, http.StatusOK, writer.Code)

        // marshal expected value using expected data 
        // to compare with the actual/ response body
        want, err := json.Marshal(UserToUserResponse(*users[0]))
        assert.NoError(t, err)

        // make sure expected body value match with the actual/ response body value
        assert.Equal(t, want, writer.Body.Bytes())
    })

    // EXPECT FAIL error param id, will return 400/ status bad request
    // we simulate this by sending empty param id
    t.Run("EXPECT FAIL error param id", func(t *testing.T){
        // prepare request/ writer/ context 
        writer, context := NewTestRecordWriter()

        // create new request with 'DELETE' method
        context.Request, _ = http.NewRequest("DELETE", "/", nil)
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed to test it behaviour
        handler.UserDeleteHandler(context)

        // make sure the expected status code (400/ bad request)
        // are match with the response code
        assert.Equal(t, http.StatusBadRequest, writer.Code)
    })

    // EXPECT FAIL error data not found, will return 500/ status internal server error
    // we simulate this by sending param id that not available on users data slice
    t.Run("EXPECT FAIL data not found", func(t *testing.T){
        // prepare request/ writer/ context 
        writer, context := NewTestRecordWriter()

        // prepare request param id. we will simulate to delete 
        // user id with id = 7 to get internal server error
        context.Params = gin.Params{
            {Key: "id", Value:"7"},
        }

        // create new request with 'DELETE' method
        context.Request, _ = http.NewRequest("DELETE", "/", nil)
        context.Request.Header.Add("content-type", "application/json")

        // actual method executed to test it behaviour
        handler.UserDeleteHandler(context)

        // make sure the expected status code (500/ internal server error)
        // are match with the response code
        assert.Equal(t, http.StatusInternalServerError, writer.Code)
    })
}
