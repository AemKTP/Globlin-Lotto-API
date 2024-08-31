package res

import (
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func GETAllLotteryResults(c *gin.Context) {
	var lotterys []models.GetLottery

	query := `	SELECT lottery.lotteryID, lottery.lotteryNumber
				FROM lottery
				LEFT JOIN winner ON lottery.lotteryID = winner.lotteryID
				WHERE winner.lotteryID IS NOT NULL
				ORDER BY winner.winnerID ASC
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error1"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lottery models.GetLottery
		if err := rows.Scan(&lottery.LotteryID, &lottery.LotteryNumber); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
			return
		}
		lotterys = append(lotterys, lottery)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error3"})
		return
	}
	c.JSON(http.StatusOK, lotterys)
}
