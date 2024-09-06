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

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userName", claims.UserName)
		c.Next()
	}

}

// GetUserIDFromContext เป็นฟังก์ชันที่ดึง userID จาก context และแปลงเป็น int64
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	// ดึง userID จาก context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return 0, gin.Error{
			Err:  gin.Error{Err: nil, Meta: "User ID not found in token"},
			Type: gin.ErrorTypePublic,
		}
	}

	// แปลง userID เป็น int64
	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID type assertion failed"})
		return 0, gin.Error{
			Err:  gin.Error{Err: nil, Meta: "User ID type assertion failed"},
			Type: gin.ErrorTypePublic,
		}
	}

	return userIDInt, nil
}
