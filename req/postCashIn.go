package req

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func CashIn(c *gin.Context) {
	var lottery models.GetLottery

	// ดึง userID จาก path parameter
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam) // แปลง userID เป็น int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
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
	err = db.DB.QueryRow(queryLottery, lottery.LotteryNumber).Scan(&lotteryID)
	if err != nil {
		log.Printf("Error finding lotteryID: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Lottery number not found"})
		return
	}

	// ยอดเงินในบัญชี
	var userBalance int
	queryUser := `SELECT userBalance FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUser, userID).Scan(&userBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			// ไม่พบผู้ใช้ในระบบ
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			// ข้อผิดพลาดในการค้นหา
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error finding user"})
		}
		return
	}

	// ตรวจสอบว่า user คนนี้มี lottery ใบนี้จริงมั้ย ผลลัพธ์ที่มี timestamp ล่าสุด
	checklotteryquery := `	SELECT userID, lotteryID
							FROM payment
							WHERE payment.lotteryID = ?
							AND transactionType = 1
							ORDER BY timestamp DESC
							LIMIT 1`

	var foundLotteryID int
	var ownerUserID int
	err = db.DB.QueryRow(checklotteryquery, lotteryID).Scan(&ownerUserID, &foundLotteryID)
	if err != nil {
		if err == sql.ErrNoRows {
			// ไม่พบข้อมูลที่ตรงกับเงื่อนไข
			c.JSON(http.StatusNotFound, gin.H{"error": "Lottery not found for this user"})
			return
		} else {
			// เกิดข้อผิดพลาดอื่น ๆ
			log.Printf("Error querying database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
			return
		}
	}

	// ตรวจสอบว่าผู้ที่ส่งคำขอเป็นเจ้าของลอตเตอรี่หรือไม่
	if ownerUserID != userID {
		// ถ้า userID ไม่ตรงกัน แสดงว่าไม่ใช่เจ้าของลอตเตอรี่ใบนี้
		c.JSON(http.StatusUnauthorized, gin.H{"error": "This lottery does not belong to the user"})
		return
	}

	// ค้นหาผลลัพธ์ของหวยในตาราง winner
	var winnerID int
	queryCheckLottery := `SELECT winnerID FROM winner WHERE lotteryID = ?`
	err = db.DB.QueryRow(queryCheckLottery, lotteryID).Scan(&winnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			// ไม่พบผลลัพธ์ของหวย
			c.JSON(http.StatusNotFound, gin.H{"error": "Sorry, you didn't win the prize."})
		} else {
			// ข้อผิดพลาดในการค้นหา
			log.Printf("Error finding lottery result: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error finding lottery result"})
		}
		return
	}
	// ตรวจสอบและดึงราคาของรางวัลจากตาราง settings
	var winnerPrize int
	switch winnerID {
	case 1:
		err = db.DB.QueryRow(`SELECT winnerPrize1 FROM settings`).Scan(&winnerPrize)
	case 2:
		err = db.DB.QueryRow(`SELECT winnerPrize2 FROM settings`).Scan(&winnerPrize)
	case 3:
		err = db.DB.QueryRow(`SELECT winnerPrize3 FROM settings`).Scan(&winnerPrize)
	case 4:
		err = db.DB.QueryRow(`SELECT winnerPrize4 FROM settings`).Scan(&winnerPrize)
	case 5:
		err = db.DB.QueryRow(`SELECT winnerPrize5 FROM settings`).Scan(&winnerPrize)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid winnerID"})
		return
	}

	if err != nil {
		log.Printf("Error finding settings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error Select lotteryPrice"})
		return
	}

	// Update เงินให้ user ตามที่ชนะ
	_, err = db.DB.Exec(`UPDATE users SET userBalance = userBalance + ? WHERE userID = ?`, winnerPrize, userID)
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

	// อัพเดท transactionType และเวลา
	_, err = db.DB.Exec(`UPDATE payment SET transactionType = ?, timestamp = ? WHERE lotteryID = ? AND userID = ? AND transactionType = 1`,
		2, timeset, lotteryID, userID)
	if err != nil {
		log.Printf("Error updating payment data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error updating payment data"})
		return
	}

	// Select userBalance ใหม่หลังจากการซื้อ
	var checkUserBalance int
	queryUserBalance := `SELECT userBalance FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUserBalance, userID).Scan(&checkUserBalance)
	if err != nil {
		log.Printf("Error finding user balance: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error selecting user balance"})
		return
	}

	// ส่งข้อความตอบกลับ
	c.JSON(http.StatusOK, gin.H{"message": "Cash-in successful", "UserID": userID, "LotteryNumber": lottery.LotteryNumber, "Winner": winnerID, "Remaining Wallet": checkUserBalance, "WinnerPrize": winnerPrize})
}
