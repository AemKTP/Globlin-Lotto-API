package res

import (
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func GetlotteriesSearch(c *gin.Context) {
	// รับค่า lotteryNumber จาก URL parameter
	lotteryNumberParam := c.Param("lotterynumber")

	// ค้นหาหมายเลขหวยที่มีการจับคู่บางส่วน
	queryLottery := `SELECT lotteryID, lotteryNumber FROM lottery WHERE lotteryNumber LIKE ?`
	rows, err := db.DB.Query(queryLottery, "%"+lotteryNumberParam+"%")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	// ดึงผลลัพธ์และสร้าง slice ของผลลัพธ์
	var lotteries []models.GetLottery
	for rows.Next() {
		var lottery models.GetLottery
		if err := rows.Scan(&lottery.LotteryID, &lottery.LotteryNumber); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		lotteries = append(lotteries, lottery)
	}

	// ตรวจสอบว่ามีข้อมูลหรือไม่
	if len(lotteries) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No lotteries found"})
		return
	}

	// ส่งข้อมูลตอบกลับ
	c.JSON(http.StatusOK, gin.H{"lotteries": lotteries})
}
