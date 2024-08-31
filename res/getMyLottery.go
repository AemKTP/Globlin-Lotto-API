package res

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func GetMyLottery(c *gin.Context) {
	var lotterys []models.GetLottery
	userIDParam := c.Param("userID")
	userID, err := strconv.Atoi(userIDParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}
	query := `SELECT 	lottery.lotteryID, lottery.lotteryNumber
			  FROM 		lottery
			  LEFT 		JOIN payment ON lottery.lotteryID = payment.lotteryID
			  WHERE 	payment.UserID = ?
			  AND		payment.transactionType = 1
			  ORDER BY 	lottery.lotteryID ASC`

	rows, err := db.DB.Query(query, userID)
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

	c.JSON(http.StatusOK, gin.H{"userID": userID, "lottery": lotterys})
}
