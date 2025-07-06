package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func actorAdminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateActor)
	r.GET("", controllers.GetAllActors)
	r.DELETE("/:id", controllers.DeleteActor)
}
