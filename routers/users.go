package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func userRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.GET("", controllers.GetAllUsers)
	r.DELETE("/:id", controllers.DeleteUserByID)
}


func profileRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.GET("", controllers.GetProfile)
	r.PATCH("", controllers.UpdateProfile)
}
