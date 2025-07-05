package routers

import (
	"be-tickitz/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CombineRouter(r *gin.Engine) {
	registerRouter(r.Group("/register"))
	loginRouter(r.Group("/login"))
	forgotPasswordRouter(r.Group("/forgot-password"))
	resetPasswordRouter(r.Group("/reset-password"))
	movieAdminRouter(r.Group("/admin/movies"))
	moviePublicRouter(r.Group("/movies"))
	genreAdminRouter(r.Group("/admin/genres"))
	genrePublicRouter(r.Group("/genres"))
	directorAdminRouter(r.Group("/admin/directors"))
  actorAdminRouter(r.Group("/admin/actors"))      
	adminPaymentMethod(r.Group("/admin/payment-method"))
	

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/docs/index.html")
	})
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
