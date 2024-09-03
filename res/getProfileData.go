package res

import (
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	var customers []models.GetDataCustomer

	// ดึง userID จาก context แทน
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	query := `SELECT userID, userName, userBalance
              FROM users
              WHERE userID = ?`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var customer models.GetDataCustomer
		if err := rows.Scan(&customer.UserID, &customer.UserName, &customer.UserBalance); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"customers": customers})
}
