package main

import "github.com/gin-gonic/gin"

func main() {
    // gin.SetMode(gin.S)
    r := gin.Default()

    r.Run(":8000")
}
