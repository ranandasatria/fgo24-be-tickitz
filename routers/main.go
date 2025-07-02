package routers

import "github.com/gin-gonic/gin"

func CombineRouter(r *gin.Engine) {
	registerRouter(r.Group("/register"))
	loginRouter(r.Group("/login"))
	forgotPasswordRouter(r.Group("/forgot-password"))
	resetPasswordRouter(r.Group("/reset-password"))	
}