package routers

import (
	"be-tickitz/controllers"
	"be-tickitz/middlewares"

	"github.com/gin-gonic/gin"
)

func genreAdminRouter(r *gin.RouterGroup){
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.CreateGenre)
}
func genrePublicRouter(r *gin.RouterGroup){
	r.GET("", controllers.GetAllGenres)
}