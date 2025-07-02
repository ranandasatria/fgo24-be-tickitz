package controllers

import (
	"be-tickitz/models"
	"be-tickitz/utils"
	"fmt"
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

	token, err := utils.GenerateJWT("auth", user.ID, 24*time.Hour, map[string]any{
		"role": user.Role,
	})
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

	token, err := utils.GenerateJWT("reset_password", user.ID, 10*time.Minute, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to generate token",
			Errors:  err.Error(),
		})
		return
	}

	body := fmt.Sprintf("<p>Copy to this token to reset your password:</p><code>%s</code>", token)

	err = utils.SendEmail(req.Email, "Reset Your Password", body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to send email",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Reset token sent to your email",
	})
}

func ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
		return
	}

	claims, err := utils.ParseJWT(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid or expired token",
			Errors:  err.Error(),
		})
		return
	}

	purpose, ok := claims["purpose"].(string)
	if !ok || purpose != "reset_password" {
		c.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Invalid token purpose",
		})
		return
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid token payload",
		})
		return
	}

	err = models.UpdateUserPassword(int(userID), req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to update password",
			Errors:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Password updated successfully",
	})
}
