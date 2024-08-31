package res

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/gin-gonic/gin"
)

func GetCheckLotteryResult(c *gin.Context) {
	// รับค่า lotteryResult จาก URL parameter
	lotteryResultParam := c.Param("lotteryResult")

	// หา lotteryID จาก lotteryResultParam
	var lotteryID int
	queryLottery := `SELECT lotteryID FROM lottery WHERE lotteryNumber = ?`
	err := db.DB.QueryRow(queryLottery, lotteryResultParam).Scan(&lotteryID)
	if err != nil {
		if err == sql.ErrNoRows {
			// ไม่พบหมายเลขหวยในฐานข้อมูล
			c.JSON(http.StatusNotFound, gin.H{"error": "Sorry, you didn't win the prize."})
		} else {
			// ข้อผิดพลาดในการค้นหา
			log.Printf("Error finding lotteryID: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error finding lotteryID"})
		}
		return
	}

	// ค้นหาผลลัพธ์ของหวยในตาราง winner
	var winnerResult int
	queryCheckLottery := `SELECT winnerID FROM winner WHERE lotteryID = ?`
	err = db.DB.QueryRow(queryCheckLottery, lotteryID).Scan(&winnerResult)
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

	// ส่งข้อความตอบกลับพร้อมผลลัพธ์ของหวย
	c.JSON(http.StatusOK, gin.H{"lotteryID": lotteryID, "winnerResult": winnerResult, "lotteryNumber": lotteryResultParam})
}
