package main

import (
	"be-tickitz/routers"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Be-Tickitz API
// @version 1.0
// @description This is a simple movie ticketing API
// @host localhost:8888
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"Message": "Backend is running"})
	})

	routers.CombineRouter(r)

	godotenv.Load()
	r.Run(fmt.Sprintf("0.0.0.0:%s", os.Getenv("APP_PORT")))
}
