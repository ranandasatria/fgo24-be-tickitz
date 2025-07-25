package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func TransactionRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateTransaction)
  r.GET("", controllers.GetMyTransactions)
}

func TransactionAdminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.GET("", controllers.GetAllTransactions)
}

func CheckSeatsRouter(r *gin.RouterGroup){
	r.GET("", controllers.CheckSeatAvailability)
}