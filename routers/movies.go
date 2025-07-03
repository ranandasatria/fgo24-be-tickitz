package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func movieAdminRouter(r *gin.RouterGroup) {
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateMovie)
	r.POST("/add-genre", controllers.AddGenretoMovie)
	r.DELETE("/:id", controllers.DeleteMovie)
}

func moviePublicRouter(r *gin.RouterGroup) {
	r.GET("", controllers.GetAllMovies)
	r.GET("/:id", controllers.GetMovieByID)
	r.GET("/now-showing", controllers.GetNowShowing)
	r.GET("/upcoming", controllers.GetUpcoming)
}
