package middleware

import (
	"net/http"
	"strings"

	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware ตรวจสอบสิทธิ์ของ Admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// สมมุติว่าคุณมีฟังก์ชันตรวจสอบสิทธิ์ของ Admin
		userType, exists := c.Get("userType")
		if !exists || userType != 1 { // Admin เช็คว่า userType เป็น 1 มั้ย
			c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden: Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserMiddleware ตรวจสอบสิทธิ์ของผู้ใช้ทั่วไป
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// สมมุติว่าคุณมีฟังก์ชันตรวจสอบสิทธิ์ของผู้ใช้
		userType, exists := c.Get("userType")
		if !exists || userType != 0 { // User เช็คว่า userType เป็น 0 มั้ย
			c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden: Users only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// PublicMiddleware ไม่มีการตรวจสอบสิทธิ์
func PublicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ไม่มีการตรวจสอบสิทธิ์ใน PublicMiddleware
		c.Next()
	}
}

func AuthenticateJWT(JWTKey []byte) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// ตรวจสอบว่า header มี "Bearer" หรือไม่
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &models.GetLoginCustomer{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Next()
	}
}
