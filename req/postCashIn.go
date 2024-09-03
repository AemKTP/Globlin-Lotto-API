package req

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func CashIn(c *gin.Context) {
	var lottery models.GetLottery

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

	// Bind JSON จาก request body ไปเก็บยังโครงสร้าง lottery
	if err := c.ShouldBindJSON(&lottery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// หา lotteryID จาก lotteryNumber ที่ส่งมาจาก body
	var lotteryID int
	queryLottery := `SELECT lotteryID FROM lottery WHERE lotteryNumber = ?`
	err := db.DB.QueryRow(queryLottery, lottery.LotteryNumber).Scan(&lotteryID)
	if err != nil {
		log.Printf("Error finding lotteryID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Lottery number not found"})
		return
	}

	// ยอดเงินในบัญชี
	var userBalance int
	queryUser := `SELECT userBalance FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUser, userIDInt).Scan(&userBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error finding user"})
		}
		return
	}

	// ตรวจสอบว่าผู้ใช้คนนี้มี lottery ใบนี้จริงหรือไม่
	checklotteryquery := `SELECT userID FROM payment WHERE lotteryID = ? AND transactionType = 1 ORDER BY timestamp DESC LIMIT 1`
	var ownerUserID int
	err = db.DB.QueryRow(checklotteryquery, lotteryID).Scan(&ownerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Lottery not found for this user"})
			return
		}
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error querying lottery"})
		return
	}

	if ownerUserID != int(userIDInt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This lottery does not belong to the user"})
		return
	}

	// ค้นหาผลลัพธ์ของหวยในตาราง winner
	var winnerID int
	queryCheckLottery := `SELECT winnerID FROM winner WHERE lotteryID = ?`
	err = db.DB.QueryRow(queryCheckLottery, lotteryID).Scan(&winnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Sorry, you didn't win the prize."})
		} else {
			log.Printf("Error finding lottery result: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error finding lottery result"})
		}
		return
	}

	// ตรวจสอบและดึงราคาของรางวัลจากตาราง settings
	var winnerPrize int
	prizeQuery := `SELECT CASE ? 
		WHEN 1 THEN winnerPrize1
		WHEN 2 THEN winnerPrize2
		WHEN 3 THEN winnerPrize3
		WHEN 4 THEN winnerPrize4
		WHEN 5 THEN winnerPrize5
	END FROM settings`
	err = db.DB.QueryRow(prizeQuery, winnerID).Scan(&winnerPrize)
	if err != nil {
		log.Printf("Error finding settings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error selecting winner prize"})
		return
	}

	// Update เงินให้ user ตามที่ชนะ
	_, err = db.DB.Exec(`UPDATE users SET userBalance = userBalance + ? WHERE userID = ?`, winnerPrize, userIDInt)
	if err != nil {
		log.Printf("Error updating user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error updating user balance"})
		return
	}

	// ใช้เวลาในปัจจุบันจาก Go แทนการ query ฐานข้อมูล
	now := time.Now()

	// อัพเดท transactionType และเวลา
	_, err = db.DB.Exec(`UPDATE payment SET transactionType = ?, timestamp = ? WHERE lotteryID = ? AND userID = ? AND transactionType = 1`,
		2, now, lotteryID, userIDInt)
	if err != nil {
		log.Printf("Error updating payment data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error updating payment data"})
		return
	}

	// Select userBalance ใหม่หลังจากการซื้อ
	var checkUserBalance int
	queryUserBalance := `SELECT userBalance FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUserBalance, userIDInt).Scan(&checkUserBalance)
	if err != nil {
		log.Printf("Error finding user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error selecting user balance"})
		return
	}

	// ส่งข้อความตอบกลับ
	c.JSON(http.StatusOK, gin.H{
		"message":          "Cash-in successful",
		"UserID":           userIDInt,
		"LotteryNumber":    lottery.LotteryNumber,
		"Winner":           winnerID,
		"Remaining Wallet": checkUserBalance,
		"WinnerPrize":      winnerPrize,
	})
}
