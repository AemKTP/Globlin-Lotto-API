package req

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

// อนาคตจะเป็นการ Reset System
func Random(c *gin.Context) {
	const numberOfNumbers = 100
	var lottery models.GetLottery

	// Bind the incoming JSON to the lottery model
	if err := c.ShouldBindJSON(&lottery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}
	lotteryNumbers := make(map[string]struct{})

	rand.Seed(time.Now().UnixNano())

	for len(lotteryNumbers) < numberOfNumbers {
		num := fmt.Sprintf("%06d", rand.Intn(1000000)) // 0-999999
		lotteryNumbers[num] = struct{}{}
	}

	numbersSlice := make([]string, 0, len(lotteryNumbers))
	for num := range lotteryNumbers {
		numbersSlice = append(numbersSlice, num)
	}

	for _, num := range numbersSlice {
		_, err := db.DB.Exec("INSERT INTO lottery (lotteryNumber) VALUES (?)", num)
		if err != nil {
			log.Printf("Error inserting user into database: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Add Lottery Successfull", "lotteryNumbers": numbersSlice})
}
