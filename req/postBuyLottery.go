package req

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/middleware"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func BuyLottery(c *gin.Context) {
	var lotterys []models.GetLottery

	// ดึง userID จาก context โดยใช้ฟังก์ชัน GetUserIDFromContext
	userIDInt, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		// ถ้ามี error ก็จะทำการ return error จาก GetUserIDFromContext
		return
	}

	// Bind JSON จาก request body ไปเก็บยังโครงสร้าง lotterys (แบบ array)
	if err := c.ShouldBindJSON(&lotterys); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// ตรวจสอบว่ามีข้อมูลลอตเตอรี่ในรายการหรือไม่
	if len(lotterys) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one lottery number is required"})
		return
	}

	// Select userBalance จาก userID
	var userBalance int
	queryUser := `SELECT userBalance FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUser, userIDInt).Scan(&userBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// หาราคา lotteryPrice จาก table settings
	var lotteryPrice int
	querySetting := "SELECT lotteryPrice FROM settings"
	err = db.DB.QueryRow(querySetting).Scan(&lotteryPrice)
	if err != nil {
		log.Printf("Error finding settings: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error selecting lotteryPrice"})
		return
	}

	// คำนวณยอดเงินรวมที่ต้องใช้ในการซื้อ
	totalPrice := len(lotterys) * lotteryPrice

	// เทียบว่าเงินในบัญชีมีพอสำหรับการซื้อทั้งหมดหรือไม่
	if userBalance < totalPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance for buying all lotteries"})
		return
	}

	// เตรียมข้อมูลสำหรับการซื้อและอัพเดทข้อมูล
	timestamp := time.Now()
	tx, err := db.DB.Begin() // เริ่ม transaction
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error starting transaction"})
		return
	}

	// สร้าง map เพื่อเก็บ lotteryID และข้อมูลที่ซื้อสำเร็จ
	lotteryIDMap := make(map[string]int)
	lotteryNumbersOutstock := make([]string, 0)
	purchasedLotteryNumbers := make([]string, 0)

	// หา lotteryID ทั้งหมดในครั้งเดียว
	lotteryNumbers := make([]string, len(lotterys))
	for i, lottery := range lotterys {
		lotteryNumbers[i] = lottery.LotteryNumber
	}

	// สร้าง placeholders สำหรับคำสั่ง SQL
	placeholders := strings.Repeat("?,", len(lotteryNumbers)-1) + "?"

	// ใช้ placeholders สำหรับคำสั่ง SQL
	queryLotteryIDs := `SELECT lotteryNumber, lotteryID FROM lottery WHERE lotteryNumber IN (` + placeholders + `)`
	rows, err := tx.Query(queryLotteryIDs, toInterfaceSlice(lotteryNumbers)...)
	if err != nil {
		tx.Rollback()
		log.Printf("Error finding lotteryIDs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error finding lotteryIDs"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lotteryNumber string
		var lotteryID int
		if err := rows.Scan(&lotteryNumber, &lotteryID); err != nil {
			tx.Rollback()
			log.Printf("Error scanning lotteryID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error scanning lotteryID"})
			return
		}
		lotteryIDMap[lotteryNumber] = lotteryID
	}

	// สร้าง slice สำหรับเก็บข้อมูลที่จะ insert
	var valueStrings []string
	var valueArgs []interface{}

	// ตรวจสอบว่า lotteryID นี้มีคนซื้อไปรึยัง
	queryCheckLottery := `SELECT lotteryID FROM payment WHERE lotteryID = ?`
	for _, lottery := range lotterys {
		lotteryID, exists := lotteryIDMap[lottery.LotteryNumber]
		if !exists {
			lotteryNumbersOutstock = append(lotteryNumbersOutstock, lottery.LotteryNumber)
			continue
		}

		var checkLottery int
		err = tx.QueryRow(queryCheckLottery, lotteryID).Scan(&checkLottery)
		if err == nil {
			lotteryNumbersOutstock = append(lotteryNumbersOutstock, lottery.LotteryNumber)
			continue
		} else if err != sql.ErrNoRows {
			tx.Rollback()
			log.Printf("Error checking lotteryID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// เพิ่มข้อมูลลงใน slice สำหรับ Batch Insert
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, userIDInt, lotteryID, 1, timestamp)

		// เพิ่ม lotteryNumber ลงใน slice ที่เก็บรายการที่ซื้อสำเร็จ
		purchasedLotteryNumbers = append(purchasedLotteryNumbers, lottery.LotteryNumber)
	}

	// ทำการ Batch Insert
	if len(valueStrings) > 0 {
		stmt := `INSERT INTO payment (userID, lotteryID, transactionType, timestamp) VALUES ` + strings.Join(valueStrings, ",")
		_, err := tx.Exec(stmt, valueArgs...)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting data into payment: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error inserting into payment"})
			return
		}
	}

	// อัพเดท userBalance หลังจากซื้อ
	if len(purchasedLotteryNumbers) > 0 {
		_, err = tx.Exec(`UPDATE users SET userBalance = userBalance - ? WHERE userID = ?`, len(purchasedLotteryNumbers)*lotteryPrice, userIDInt)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating user balance: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error updating userBalance"})
			return
		}
	}

	// ถ้าทุกอย่างผ่านไปได้ดี ก็ commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error committing transaction"})
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
	response := gin.H{
		"Purchased Lottery Numbers": purchasedLotteryNumbers,
		"Remaining Wallet":          checkUserBalance,
		"message":                   "Buy Lottery process completed",
	}
	if len(lotteryNumbersOutstock) > 0 {
		response["Lottery Numbers OutStock"] = lotteryNumbersOutstock
	}
	c.JSON(http.StatusOK, response)
}

// Helper function to convert a slice of strings to a slice of interface{}
func toInterfaceSlice(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
