package controllers

import (
	"be-tickitz/models"
	"be-tickitz/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
		return
	}

	createdUser, err := models.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to create user",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User created",
		Results: createdUser,
	})
}

func Login(ctx *gin.Context) {
	godotenv.Load()
	secretKey := os.Getenv("APP_SECRET")

	form := struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	user, err := models.FindOneUserByEmail(form.Email)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	if err := utils.CompareHash(user.Password, form.Password); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Wrong email or password",
		})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	token, _ := generateToken.SignedString([]byte(secretKey))

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Login success",
		Results: map[string]string{
			"token": token,
		},
	})
}
