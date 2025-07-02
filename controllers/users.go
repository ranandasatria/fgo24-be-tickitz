package controllers

import (
	"be-tickitz/models"
	"be-tickitz/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	token, err := utils.GenerateJWT("auth", user.ID, 24*time.Hour)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to generate token",
			Errors:  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Login success",
		Results: map[string]string{
			"token": token,
		},
	})
}

func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid email format",
		})
		return
	}

	user, err := models.FindOneUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.Response{
			Success: false,
			Message: "Email not found",
		})
		return
	}

	token, err := utils.GenerateJWT("reset_password", user.ID, 10*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to generate token",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Reset token created",
		Results: gin.H{"resetToken": token},
	})
}
