/*
   package account
   handler.go
       handler/ controller layer for account package. outer layer that interact with client
       which processing user request and response
*/
package account

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// accountHandler is type wrapper for AccountService interface
type accountHandler struct {
    Service AccountService
}

// NewAccountHandler is new instance for accountHandler
func NewAccountHandler(svc AccountService) accountHandler{
    return accountHandler{Service: svc}
}

// UserCreate method will process request to insert new 'User' data and
// response with the created data back to the user (if no error found)
func (h *accountHandler) UserCreateHandler(c *gin.Context) {
    // get request data from context that containing 'User' model information
    // and bind it to a variable matching the requested data
    var u User
    

    // if request data binding error than return 400/ bad request
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(
            http.StatusBadRequest, 
            gin.H{
                "error": fmt.Sprintf("bad request: %v\n", err),
            },
        )

        // exit process
        return
    }

    // send data to service layer to further process (create record)
    user, err := h.Service.Create(u)

    // if error occur while trying to save the data, return 500/ internal server error
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusInternalServerError,
            gin.H{
                "error": fmt.Sprintf("internal server error: %v\n", err),
            },
        )

        // exit process
        return
    }

    //  if no error found, send 200/ status ok as well as the 'UserResponse' data
    c.JSON(
        http.StatusOK,
        user,
    )
}

// UserGetHandler is method to process request to get user data 
func (h *accountHandler) UserGetHandler(c *gin.Context) {
    id := c.Param("id")
    uid, err := strconv.Atoi(id)
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusBadRequest,
            gin.H{
                "error": fmt.Sprintf("bad request: %v\n", err),
            },
        )
        return
    }

    user, err := h.Service.Get(uid)
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusInternalServerError,
            gin.H{
                "error": fmt.Sprintf("internal server error: %v\n", err),
            },
        )
        return
    }

    c.JSON(http.StatusOK, user)
}

// UserGetsHandler is method to process request to get all user data 
func (h *accountHandler) UserGetsHandler(c *gin.Context) {
    users, err := h.Service.Gets()
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusInternalServerError,
            gin.H{
                "error": fmt.Sprintf("internal server error: %v\n", err),
            },
        )
        return
    }

    // no error occur then send status ok and users data
    c.JSON(
        http.StatusOK, 
        users,
    )
}

// UserUpdateHandler method will process request to update 'User' data and
// response with the updated data
func (h *accountHandler) UserUpdateHandler(c *gin.Context) {
    // get request parameter for 'id'
    id := c.Param("id")
    uid, err := strconv.Atoi(id)

    // if error found, respnse with 400(bad request) and exit the process
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusBadRequest,
            gin.H{
                "error": fmt.Sprintf("bad request: %v\n", err),
            },
        )

        // exit process
        return
    }
    // get request data from context that containing 'User' model information
    // and bind it to a variable matching the requested data
    var u User

    // if request data binding error than return 400/ bad request
    if err := c.ShouldBindJSON(&u); err != nil {
        c.JSON(
            http.StatusBadRequest, 
            gin.H{
                "error": fmt.Sprintf("bad request: %v\n", err),
            },
        )

        // exit process
        return
    }

    // send data to service layer to further process (update record)
    user, err := h.Service.Update(uid, u)

    // if error occur while trying to save the data, return 500/ internal server error
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusInternalServerError,
            gin.H{
                "error": fmt.Sprintf("internal server error: %v\n", err),
            },
        )

        // exit process
        return
    }

    //  if no error found, send 200/ status ok as well as the 'UserResponse' data
    c.JSON(
        http.StatusOK,
        user,
    )
}

// UserDeleteHandler method will process request to delete 'User' data and
// response with the updated data
func (h *accountHandler) UserDeleteHandler(c *gin.Context) {
    // get request parameter for 'id'
    id := c.Param("id")
    uid, err := strconv.Atoi(id)

    // if error found, respnse with 400(bad request) and exit the process
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusBadRequest,
            gin.H{
                "error": fmt.Sprintf("bad request: %v\n", err),
            },
        )

        // exit process
        return
    }

    // send data to service layer to further process (delete record)
    user, err := h.Service.Delete(uid)

    // if error occur while trying to save the data, return 500/ internal server error
    if err != nil {
        c.AbortWithStatusJSON(
            http.StatusInternalServerError,
            gin.H{
                "error": fmt.Sprintf("internal server error: %v\n", err),
            },
        )

        // exit process
        return
    }

    //  if no error found, send 200/ status ok as well as the 'UserResponse' data
    c.JSON(
        http.StatusOK,
        user,
    )
}
