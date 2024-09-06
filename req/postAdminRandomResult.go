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
	"github.com/AemKTP/Globlin-Lotto-API/middleware"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func RandomResult(c *gin.Context) {
	const NumbersResult = 5

	// ดึง userID จาก context โดยใช้ฟังก์ชัน GetUserIDFromContext
	userIDInt, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		// ถ้ามี error ก็จะทำการ return error จาก GetUserIDFromContext
		return
	}

	var userType int
	queryUser := `SELECT userType FROM users WHERE userID = ?`
	err = db.DB.QueryRow(queryUser, userIDInt).Scan(&userType)
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

	// ลบข้อมูลจากตาราง winner
	_, err = tx.Exec("TRUNCATE TABLE winner")
	if err != nil {
		log.Printf("Error deleting from winner table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error deleting winner data"})
		return
	}

	// สุ่มเลือก lottery 5 รายการ
	rows, err := tx.Query("SELECT lotteryID FROM lottery")
	if err != nil {
		log.Printf("Error selecting lottery IDs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error selecting lottery IDs"})
		return
	}
	defer rows.Close()

	var lotteryIDs []int
	for rows.Next() {
		var lotteryID int
		if err := rows.Scan(&lotteryID); err != nil {
			log.Printf("Error scanning lottery ID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error scanning lottery ID"})
			return
		}
		lotteryIDs = append(lotteryIDs, lotteryID)
	}

	if len(lotteryIDs) < NumbersResult {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough lottery entries in the database"})
		return
	}

	// สุ่มเลือก lottery 5 รายการ
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(lotteryIDs), func(i, j int) { lotteryIDs[i], lotteryIDs[j] = lotteryIDs[j], lotteryIDs[i] })

	selectedLotteryIDs := lotteryIDs[:NumbersResult]

	// เตรียมข้อมูลสำหรับ Batch Insert
	valueStrings := make([]string, 0, len(selectedLotteryIDs))
	valueArgs := make([]interface{}, 0, len(selectedLotteryIDs))

	for _, id := range selectedLotteryIDs {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, id)
	}

	// Batch Insert
	stmt := fmt.Sprintf("INSERT INTO winner (lotteryID) VALUES %s", strings.Join(valueStrings, ","))
	_, err = tx.Exec(stmt, valueArgs...)
	if err != nil {
		log.Printf("Error inserting into winner table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error inserting winner data"})
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error committing transaction"})
		return
	}

	// ดึงข้อมูลจากตาราง winner
	var winners []models.Winner
	rows, err = db.DB.Query("SELECT winnerID, lotteryID FROM winner ORDER BY winnerID ASC")
	if err != nil {
		log.Printf("Error selecting winners: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error selecting winners"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var winner models.Winner
		if err := rows.Scan(&winner.WinnerID, &winner.LotteryID); err != nil {
			log.Printf("Error scanning winner data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error scanning winner data"})
		}
		winners = append(winners, winner)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Random lottery results successful", "winners": winners})
}
