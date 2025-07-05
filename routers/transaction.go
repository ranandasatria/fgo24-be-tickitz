package routers

import (
  "be-tickitz/controllers"
  "be-tickitz/middlewares"
  "github.com/gin-gonic/gin"
)

func TransactionRouter(r *gin.RouterGroup) {
  r.Use(middlewares.VerifyToken())
  r.POST("", controllers.CreateTransaction)
}
