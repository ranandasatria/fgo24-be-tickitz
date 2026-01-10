package main

import (
	"be-tickitz/routers"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Be-Tickitz API
// @version 1.0
// @description This is a simple movie ticketing API
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// if os.Getenv("ENV") == "local" {
	// 	godotenv.Load()
	// }
	_ = godotenv.Load()
	r := gin.New()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"Message": "Backend is running"})
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	routers.CombineRouter(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	r.Run("0.0.0.0:" + port)

}
