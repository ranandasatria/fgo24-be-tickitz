package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func directorAdminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateDirector)
	r.GET("", controllers.GetAllDirectors)
	r.DELETE("/:id", controllers.DeleteDirector)
}
