package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.RouterGroup) {
	r.POST("", controllers.Register)
}

func loginRouter(r *gin.RouterGroup) {
	r.POST("", controllers.Login)
}

func forgotPasswordRouter(r *gin.RouterGroup) {
	r.POST("", controllers.ForgotPassword)
}
