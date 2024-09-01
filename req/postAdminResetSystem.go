package req

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

// อนาคตจะเป็นการ Reset System
func ResetSystem(c *gin.Context) {
	const numberOfNumbers = 100
	var user models.GetCustomer

	// รับค่า userID จาก URL parameter และแปลงเป็น int
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	// ถ้า Type ไม่ = 1 ก็ให้ return กลับไป
	var userType int
	queryLottery := `SELECT userID FROM users WHERE userID = ? AND userType = 1`
	err = db.DB.QueryRow(queryLottery, userID).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Delete Table
	tables := []string{"payment", "winner", "users", "lottery"}
	for _, table := range tables {
		_, err = db.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Printf("Error deleting from %s table: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error deleting from %s table", table)})
			return
		}

		// Reset AUTO_INCREMENT
		_, err = db.DB.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table))
		if err != nil {
			log.Printf("Error resetting AUTO_INCREMENT for %s table: %v", table, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error resetting AUTO_INCREMENT for %s table", table)})
			return
		}
	}

	// Create Admin
	user.UserName = "goblin123"
	user.Pwd = "$2a$10$/DCkF0KFmQuUkpQwS7i6o./pHZSKgazJa0GlcPtuwMhaxlEERsiv."
	user.UserType = 1

	// บันทึกข้อมูลลงในฐานข้อมูล
	_, err = db.DB.Exec("INSERT INTO users (userName, userPwd, userType) VALUES (?, ?, ?)", user.UserName, user.Pwd, user.UserType)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// log.Printf("userName: %s, userPwd: %s, userType: %d", user.UserName, user.Pwd, user.UserType)

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

	for _, num := range numbersSlice {
		_, err := db.DB.Exec("INSERT INTO lottery (lotteryNumber) VALUES (?)", num)
		if err != nil {
			log.Printf("Error inserting user into database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset System Successfull", "lotteryNumbers": numbersSlice})
}
