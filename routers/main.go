package routers

import "github.com/gin-gonic/gin"

func CombineRouter(r *gin.Engine) {
	registerRouter(r.Group("/register"))
}