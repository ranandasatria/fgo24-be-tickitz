package middlewares

import (
	"be-tickitz/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Missing token",
			})
			ctx.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.Response{
				Success: false,
				Message: "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("user", claims)
		ctx.Next()
	}
}
