package main

import (
	"log"
	"pgxtest/account"

	"github.com/gin-gonic/gin"
)

func main() {
    // separate the code from the 'main' function so we can test it.
    // all code that available in main function were not testable
    Run()
}

func Run() {
    // prepare gin
    // uncomment below mode if want to get back to default debug mode
    gin.SetMode(gin.ReleaseMode)

    // gin with default setup
    r := gin.New()
    // r.Use(gin.Logger())
    r.Use(gin.Recovery())

    // prepare database, remember to create "golangtest" database first
    dbPool, _, err := account.NewDBPool(account.DatabaseConfig{
        Username : "golang",
        Password : "golang",
        Hostname : "localhost",
        Port : "5432",
        DBName : "golangtest",
    })
    
    defer dbPool.Close()

    // log for error if error occur while connecting to the database
    if err != nil {
        log.Fatalf("unexpected error while tried to connect to database: %v\n", err)
    }

    // datastore := account.NewDatastore(dbPool)
    accDB := account.NewDatabase(dbPool)
    accService := account.NewAccountService(accDB)
    accAPI := account.NewAccountHandler(accService)

    // prepare router
    // main group api endpoint url : http://domain.com/v1
    v1 := r.Group("/v1")

    // account app group api endpoint : http://domainname.com/v1/account
    accRouter := v1.Group("/account")
    accRouter.POST("/", accAPI.UserCreateHandler)
    accRouter.PUT("/:id", accAPI.UserUpdateHandler)
    accRouter.DELETE("/:id", accAPI.UserDeleteHandler)
    accRouter.GET("/:id", accAPI.UserGetHandler)
    accRouter.GET("/", accAPI.UserGetsHandler)

    // run the server
    log.Fatalf("%v", r.Run(":8000"))
}
