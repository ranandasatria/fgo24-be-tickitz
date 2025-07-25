package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func adminPaymentMethod(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreatePaymentMethod)
	r.GET("", controllers.GetAllPaymentMethod)
	r.DELETE("/:id", controllers.DeletePaymentMethod)
}

func userPaymentMethod(r *gin.RouterGroup) {
	r.GET("", controllers.GetAllPaymentMethod)
}

