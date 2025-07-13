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
	userRouter(r.Group("/users"))
	profileRouter(r.Group("/profile"))
	movieAdminRouter(r.Group("/admin/movies"))
	moviePublicRouter(r.Group("/movies"))
	genreAdminRouter(r.Group("/admin/genres"))
	genrePublicRouter(r.Group("/genres"))
	directorAdminRouter(r.Group("/admin/directors"))
	directorPublicRouter(r.Group("/directors"))
	actorAdminRouter(r.Group("/admin/actors"))
	actorPublicRouter(r.Group("/actors"))
	adminPaymentMethod(r.Group("/admin/payment-method"))
	userPaymentMethod(r.Group("/payment-method"))
	TransactionRouter(r.Group("/transactions"))
	TransactionAdminRouter(r.Group("/admin/transactions"))
	CheckSeatsRouter(r.Group("/check-seats"))

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/docs/index.html")
	})
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
