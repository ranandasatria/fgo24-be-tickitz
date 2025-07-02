package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AuthRegisterLogin true "User data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /register [post]
func Register(c *gin.Context) {
	var input dto.AuthRegisterLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err.Error(),
		})
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
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

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AuthRegisterLogin true "Login data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /login [post]
func Login(ctx *gin.Context) {
	var form dto.AuthRegisterLogin

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

// ForgotPassword godoc
// @Summary Send password reset token
// @Description Send reset token to user's email if email is valid
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "User email"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /forgot-password [post]
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

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset password using valid reset token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body object{token=string,newPassword=string} true "Reset password data"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /reset-password [post]
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
