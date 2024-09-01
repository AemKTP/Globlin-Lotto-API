package req

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func BuyLottery(c *gin.Context) {
	var lotterys models.GetLottery

	// รับค่า userID จาก URL parameter และแปลงเป็น int
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	// Bind JSON จาก request body ไปเก็บยังโครงสร้าง lotterys
	if err := c.ShouldBindJSON(&lotterys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// ตรวจสอบข้อมูลที่ได้รับ
	if lotterys.LotteryNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	// หา lotteryID จาก lotteryNumber
	var lotteryID int
	queryLottery := `SELECT lotteryID FROM lottery WHERE lotteryNumber = ?`
	err = db.DB.QueryRow(queryLottery, lotterys.LotteryNumber).Scan(&lotteryID)
	if err != nil {
		log.Printf("Error finding lotteryID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Lottery number not found"})
		return
	}

	// Select userBalance จาก ID ที่ path param มา
	var userBalance int
	queryUser := `SELECT userBalance FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUser, userID).Scan(&userBalance)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// หาราคา lotteryPrice จาก table settings
	var lotteryPrice int
	querySetting := "SELECT lotteryPrice FROM settings"
	err = db.DB.QueryRow(querySetting).Scan(&lotteryPrice)
	if err != nil {
		log.Printf("Error finding settings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error Select lotteryPrice"})
		return
	}

	// เทียบว่า เงินในบัญชีมีพอสำหรับการซื้อ lottery มั้ย
	if userBalance < lotteryPrice {
		log.Printf("Insufficient balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insufficient balance"})
		return
	}

	// Update userBalance หลังจากซื้อแล้ว / ง่ายๆ คือ คำนวณเงินที่ซื้อ lottery ไป
	_, err = db.DB.Exec(`UPDATE users SET userBalance = userBalance - ? WHERE userID = ?`, lotteryPrice, userID)
	if err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error at Update userBalance"})
		return
	}

	// print select NOW() in sql
	queryNow := "SELECT NOW()"
	var nowStr string
	err = db.DB.QueryRow(queryNow).Scan(&nowStr)
	if err != nil {
		log.Printf("Error finding now: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	now, err := time.Parse("2006-01-02 15:04:05", nowStr)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	// log.Printf("Current time in Asia/Bangkok: %v", now)

	// set now to timeset
	timeset := now

	// บันทึกข้อมูลลงในฐานข้อมูล เพิ่มข้อมูลลง payment
	_, err = db.DB.Exec("INSERT INTO payment (userID, lotteryID, transactionType, timestamp) VALUES (?, ?, ?, ?)",
		userID, lotteryID, 1, timeset)
	if err != nil {
		log.Printf("Error inserting data into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// ส่งข้อความตอบกลับ
	c.JSON(http.StatusOK, gin.H{"message": "Buy Lottery successfully"})
}
