package req

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var user models.GetCustomer
	var storedUser models.GetCustomer

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	if user.UserName == "" || user.Pwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
		return
	}

	query := "SELECT userID, userName, userPwd FROM users WHERE userName = ? LIMIT 1"
	err := db.DB.QueryRow(query, user.UserName).Scan(
		&storedUser.UserID, &storedUser.UserName, &storedUser.Pwd,
	)
	if err != nil {
		log.Printf("Error finding user in database: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Pwd), []byte(user.Pwd)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	tokenString, err := generateToken(storedUser, jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": tokenString})
}

// แยก Function ออกมาให้ดูง่าย และ เป็นส่วน
func generateToken(user models.GetCustomer, jwtKey []byte) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.GetLoginCustomer{
		UserID:   user.UserID,
		UserName: user.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
