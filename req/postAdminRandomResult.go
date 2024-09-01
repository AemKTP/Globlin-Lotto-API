package req

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func RandomResult(c *gin.Context) {
	const NumbersResult = 5

	// รับค่า userID จาก URL parameter และแปลงเป็น int
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	var userType int
	queryUser := `SELECT userID FROM users WHERE userID = ? AND userType = 1`
	err = db.DB.QueryRow(queryUser, userID).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		} else {
			log.Printf("Error finding user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// ลบข้อมูลจากตาราง winner
	_, err = db.DB.Exec("TRUNCATE TABLE winner")
	if err != nil {
		log.Printf("Error deleting from winner table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error deleting winner data"})
		return
	}

	// สุ่มเลือก lottery 5 รายการ
	rows, err := db.DB.Query("SELECT lotteryID FROM lottery")
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

	// เพิ่มข้อมูลลงในตาราง winner
	for _, id := range selectedLotteryIDs {
		_, err := db.DB.Exec("INSERT INTO winner (lotteryID) VALUES (?)", id)
		if err != nil {
			log.Printf("Error inserting into winner table: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error inserting winner data"})
			return
		}
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
			return
		}
		winners = append(winners, winner)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Random lottery results successful", "winners": winners})
}
