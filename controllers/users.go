package controllers

import (
	"be-tickitz/dto"
	"be-tickitz/models"
	"be-tickitz/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	response := dto.UserResponse{
		ID:       createdUser.ID,
		Email:    createdUser.Email,
		FullName: createdUser.FullName,
		Role:     createdUser.Role,
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "User created",
		Results: response,
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

// @Summary Get all users (admin only)
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Only admin can access",
		})
		return
	}

	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  err.Error(),
		})
		return
	}

	var userList []dto.UserListResponse
	for _, u := range users {
		userList = append(userList, dto.UserListResponse{
			ID:       u.ID,
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
			Phone:    u.PhoneNumber,
			Picture:  u.ProfilePicture,
		})
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "All users",
		Results: userList,
	})
}

// @Summary Get profile detail
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	userID := int(claims["userId"].(float64))

	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed to fetch profile",
			Errors:  err.Error(),
		})
		return
	}

	response := dto.UserListResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Phone:    user.PhoneNumber,
		Picture:  user.ProfilePicture,
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Your profile",
		Results: response,
	})
}

// @Summary Delete user by ID (admin only)
// @Tags Users
// @Security BearerAuth
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 403 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /admin/users/{id} [delete]
func DeleteUserByID(c *gin.Context) {
	claims := c.MustGet("user").(jwt.MapClaims)
	role := claims["role"].(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, utils.Response{
			Success: false,
			Message: "Only admin can access",
		})
		return
	}

	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	err = models.DeleteUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, utils.Response{
				Success: false,
				Message: "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success: false,
				Message: "Failed to delete user",
				Errors:  err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: fmt.Sprintf("User with ID %d deleted", userID),
	})
}
