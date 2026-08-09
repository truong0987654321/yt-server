package middleware

import (
	"strings"
	"yt-clone-server/internal/apperrors"
	"yt-clone-server/internal/service"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"
const ContextUserEmailKey = "userEmail"

// RequireAuth kiểm tra header "Authorization: Bearer <token>".
// Dùng được cho MỌI client (Next.js, iOS app, Android app) vì JWT không phụ thuộc cookie/session.
func RequireAuth(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			appErr := apperrors.NewAuthorization("thiếu hoặc sai định dạng Authorization header")
			c.AbortWithStatusJSON(appErr.Status(), appErr)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtService.ParseAccessToken(tokenStr)
		if err != nil {
			appErr := apperrors.Parse(err)
			c.AbortWithStatusJSON(appErr.Status(), appErr)
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Next()
	}
}
