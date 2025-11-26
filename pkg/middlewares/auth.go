package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	PublicKey interface{}
}

func AuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		tokenString, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no access token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return config.PublicKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		companyID, ok := claims["cid"].(string)
		if !ok || companyID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "company id not found in token"})
			c.Abort()
			return
		}

		c.Set("company_id", companyID)
		c.Next()
	}
}

func GetCompanyID(c *gin.Context) string {
	companyID, _ := c.Get("company_id")
	if id, ok := companyID.(string); ok {
		return id
	}
	return ""
}
