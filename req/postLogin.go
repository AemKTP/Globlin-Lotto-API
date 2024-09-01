package req

import (
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var user models.GetCustomer
	var storedUser models.GetCustomer

	// Bind JSON จาก request body ไปเก็บยังโครงสร้าง user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// ตรวจสอบข้อมูลที่ได้รับ
	if user.UserName == "" || user.Pwd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password are required"})
		return
	}

	// ค้นหาผู้ใช้ในฐานข้อมูลโดยใช้ชื่อผู้ใช้
	query := "SELECT userID, userName, userPwd FROM users WHERE userName = ? LIMIT 1"
	err := db.DB.QueryRow(query, user.UserName).Scan(
		&storedUser.UserID, &storedUser.UserName, &storedUser.Pwd,
	)
	if err != nil {
		log.Printf("Error finding user in database: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	// ตรวจสอบรหัสผ่าน (compare hashed password)
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Pwd), []byte(user.Pwd)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// ส่ง response กลับไปเมื่อการล็อกอินสำเร็จ
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": storedUser})
}
