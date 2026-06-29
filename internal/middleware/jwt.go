package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/example/golang-api/pkg/ctv"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(gtx *gin.Context) {
		authHeader := gtx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			gtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"message":   "missing or invalid token",
				"code":      http.StatusUnauthorized,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			gtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"message":   "invalid or expired token",
				"code":      http.StatusUnauthorized,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			gtx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success":   false,
				"message":   "invalid token claims",
				"code":      http.StatusUnauthorized,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		ctv.Set(gtx, ctv.EmailKey, claims["email"].(string))
		gtx.Next()
	}
}
