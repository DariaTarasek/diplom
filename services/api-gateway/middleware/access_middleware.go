package middleware

import (
	"github.com/DariaTarasek/diplom/services/api-gateway/clients"
	authpb "github.com/DariaTarasek/diplom/services/api-gateway/proto/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MakeAccessMiddleware(authClient *clients.AuthClient) func(requiredPermission int32) gin.HandlerFunc {
	return func(requiredPermission int32) gin.HandlerFunc {
		return func(c *gin.Context) {
			token, err := c.Cookie("access_token")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен не найден"})
				return
			}
			_, err = authClient.Client.PermissionCheck(c.Request.Context(), &authpb.PermissionCheckRequest{
				Token:  token,
				PermId: requiredPermission,
			})

			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав"})
				return
			}
			c.Next()
		}
	}
}
