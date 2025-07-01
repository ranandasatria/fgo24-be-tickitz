package routers

import (
	"be-tickitz/controllers"

	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.RouterGroup) {
	r.POST("", controllers.Register)
}
