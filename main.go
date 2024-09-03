package main

import (
	"log"
	"os"

	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/middleware"
	"github.com/AemKTP/Globlin-Lotto-API/req"
	"github.com/AemKTP/Globlin-Lotto-API/res"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var JWTKey []byte

func init() {
	// โหลด Environment Variables จากไฟล์ .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// อ่านค่า jwtKey จาก Environment Variable
	JWTKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func main() {
	db.Init()

	r := gin.Default()

	// ใช้ middleware สำหรับเส้นทางที่ต้องการ JWT กรณีที่มี Path ที่อยากให้ใช้ JWT
	// authorized := r.Group("", middleware.AuthenticateJWT(JWTKey))
	// {
	// r.GET("/profile/:userID", res.GetProfile)
	// 	// เพิ่มเส้นทางอื่นๆ ที่ต้องการ JWT ที่นี่
	// }

	r.GET("/lotteries", res.GetLotterys)
	r.GET("/canbuylotteries", res.GetCanBuyLotteries)
	r.GET("/AllLotteryResult/", res.GETAllLotteryResults)
	r.GET("/lotteriesSearch/:lotterynumber", res.GetlotteriesSearch)
	r.GET("/CheckLotteryResult/:lotteryResult", res.GetCheckLotteryResult)
	r.GET("/profile", middleware.AuthenticateJWT(JWTKey), res.GetProfile)
	r.GET("/MyLottery", middleware.AuthenticateJWT(JWTKey), res.GetMyLottery)

	// Admin
	r.GET("/users", res.GetUsers)

	r.POST("/register", req.Register)
	r.POST("/login", req.Login)
	r.POST("/buylottery", middleware.AuthenticateJWT(JWTKey), req.BuyLottery)
	r.POST("/cashin", middleware.AuthenticateJWT(JWTKey), req.CashIn)

	// Admin
	// r.POST("/randomlotteryResult/:userID", req.RandomResult)
	// r.POST("/resetSystem/:userID", req.ResetSystem)
	r.POST("/randomlotteryResult", middleware.AuthenticateJWT(JWTKey), req.RandomResult)
	r.POST("/resetSystem", middleware.AuthenticateJWT(JWTKey), req.ResetSystem)

	r.Run(":8090")

}
