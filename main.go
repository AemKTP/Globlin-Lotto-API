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
	// เริ่มต้นเชื่อมต่อฐานข้อมูล
	db.Init()

	// ตั้งค่า Gin Router
	r := gin.Default()

	// กลุ่มเส้นทางที่ต้องการ JWT Authentication
	authorized := r.Group("", middleware.AuthenticateJWT(JWTKey))
	{
		authorized.GET("/profile", res.GetProfile)
		authorized.GET("/MyLottery", res.GetMyLottery)
		authorized.POST("/buylottery", req.BuyLottery)
		authorized.POST("/cashin", req.CashIn)

		// Admin Routes
		authorized.POST("/randomlotteryResult", req.RandomResult)
		authorized.POST("/resetSystem", req.ResetSystem)
	}

	// เส้นทางที่ไม่ต้องการ JWT Authentication
	r.GET("/lotteries", res.GetLotterys)
	r.GET("/canbuylotteries", res.GetCanBuyLotteries)
	r.GET("/AllLotteryResult/", res.GETAllLotteryResults)
	r.GET("/lotteriesSearch/:lotterynumber", res.GetlotteriesSearch)
	r.GET("/CheckLotteryResult/:lotteryResult", res.GetCheckLotteryResult)

	// เส้นทางที่เกี่ยวกับ Admin ที่อาจจะไม่ได้ใช้ JWT Authentication
	r.GET("/users", res.GetUsers)

	// เส้นทางสำหรับการลงทะเบียนและล็อกอิน
	r.POST("/register", req.Register)
	r.POST("/login", req.Login)

	// รันเซิร์ฟเวอร์บนพอร์ต 8090
	r.Run(":8090")
}
