package middleware

import (
	"net/http"
	auth_handlers "restful-service/handlers/auth"
	"restful-service/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			ctx.Abort()
			return
		}

		tokenString = strings.Split(tokenString, "Bearer ")[1]

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			ctx.Abort()
			return
		}

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			ctx.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, getSigningKey)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			ctx.Abort()
			return
		}

		_, ok := token.Claims.(*models.Claims)

		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func getSigningKey(token *jwt.Token) (interface{}, error) {
	return []byte(auth_handlers.JwtKey), nil
}
