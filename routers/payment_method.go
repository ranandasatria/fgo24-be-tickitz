package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func adminPaymentMethod(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreatePaymentMethod)
}
