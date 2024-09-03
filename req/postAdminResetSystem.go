package req

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

// อนาคตจะเป็นการ Reset System
func ResetSystem(c *gin.Context) {
	const numberOfNumbers = 100
	var user models.GetCustomer

	// ดึง userID จาก context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// แปลง userID เป็น int64
	userIDInt, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID type assertion failed"})
		return
	}

	// ตรวจสอบประเภทผู้ใช้
	var userType int
	queryUser := `SELECT userType FROM users WHERE userID = ?`
	err := db.DB.QueryRow(queryUser, userIDInt).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if userType != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// เริ่มต้น Transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error starting transaction"})
		return
	}
	defer tx.Rollback()

	// ลบข้อมูลจากตาราง
	tables := []string{"payment", "winner", "users", "lottery"}
	for _, table := range tables {
		_, err = tx.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("Error deleting from %s table: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error deleting from %s table", table)})
			return
		}

		// รีเซ็ต AUTO_INCREMENT
		_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table))
		if err != nil {
			log.Printf("Error resetting AUTO_INCREMENT for %s table: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error resetting AUTO_INCREMENT for %s table", table)})
			return
		}
	}

	// สร้างผู้ดูแลระบบ
	user.UserName = "goblin123"
	user.Pwd = "$2a$10$/DCkF0KFmQuUkpQwS7i6o./pHZSKgazJa0GlcPtuwMhaxlEERsiv."
	user.UserType = 1

	// บันทึกข้อมูลลงในฐานข้อมูล
	_, err = tx.Exec("INSERT INTO users (userName, userPwd, userType) VALUES (?, ?, ?)", user.UserName, user.Pwd, user.UserType)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Random lottery
	lotteryNumbers := make(map[string]struct{})

	rand.Seed(time.Now().UnixNano())

	for len(lotteryNumbers) < numberOfNumbers {
		num := fmt.Sprintf("%06d", rand.Intn(1000000)) // 0-999999
		lotteryNumbers[num] = struct{}{}
	}

	numbersSlice := make([]string, 0, len(lotteryNumbers))
	for num := range lotteryNumbers {
		numbersSlice = append(numbersSlice, num)
	}

	// Batch insert
	valueStrings := make([]string, 0, len(numbersSlice))
	valueArgs := make([]interface{}, 0, len(numbersSlice))
	for _, num := range numbersSlice {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, num)
	}

	stmt := fmt.Sprintf("INSERT INTO lottery (lotteryNumber) VALUES %s", strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		log.Printf("Error inserting lottery numbers into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error committing transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset System Successful", "lotteryNumbers": numbersSlice})
}
