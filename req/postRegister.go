package req

import (
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.GetCustomer

	// Bind JSON จาก request body ไปเก็บยังโครงสร้าง user
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// ตรวจสอบข้อมูลที่ได้รับ
	if user.UserName == "" || user.Pwd == "" || user.UserBalance <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	user.Pwd = string(hashedPassword)

	// บันทึกข้อมูลลงในฐานข้อมูล
	_, err = db.DB.Exec("INSERT INTO users (userName, userPwd, userBalance) VALUES (?, ?, ?)",
		user.UserName, user.Pwd, user.UserBalance)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})

}
