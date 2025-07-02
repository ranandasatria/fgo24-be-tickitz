package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func movieAdminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateMovie)
}

func moviePublicRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetAllMovies)
	r.GET("/:id", controllers.GetMovieByID)
}
