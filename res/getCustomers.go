package res

import (
	"log"
	"net/http"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/models"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.GetCustomer

	rows, err := db.DB.Query("SELECT userID, userName, userBalance FROM users")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error2"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.GetCustomer
		if err := rows.Scan(&user.UserID, &user.UserName, &user.UserBalance); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error3"})
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error4"})
		return
	}

	c.JSON(http.StatusOK, users)
}
