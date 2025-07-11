package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func actorAdminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateActor)
	r.DELETE("/:id", controllers.DeleteActor)
}

func actorPublicRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetAllActors)
}