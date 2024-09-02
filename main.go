package main

import (
	"github.com/AemKTP/Globlin-Lotto-API/db"
	"github.com/AemKTP/Globlin-Lotto-API/req"
	"github.com/AemKTP/Globlin-Lotto-API/res"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()

	r.GET("/lotteries", res.GetLotterys)                                   // show lottery
	r.GET("/canbuylotteries", res.GetCanBuyLotteries)                      // lottery ที่สามารถซื้อได้
	r.GET("/MyLottery/:userID", res.GetMyLottery)                          // MyLottery
	r.GET("/profile/:userID", res.GetProfile)                              // Myprofile
	r.GET("/lotteriesSearch/:lotterynumber", res.GetlotteriesSearch)       // Search lottery Number
	r.GET("/CheckLotteryResult/:lotteryResult", res.GetCheckLotteryResult) // Check Award Number
	r.GET("/AllLotteryResult/", res.GETAllLotteryResults)                  // lottery ที่ออกรางวัล

	// Admin
	r.GET("/users", res.GetUsers)

	r.POST("/register", req.Register)
	r.POST("/login", req.Login)
	r.POST("/buylottery/:userID", req.BuyLottery)
	r.POST("/cashin/:userID", req.CashIn)

	// Admin
	r.POST("/randomlotteryResult/:userID", req.RandomResult) //เฉพาะแอดมิน จริงๆไม่ได้เอาไปใช้ตอนโชว์ ใช้ผ่านการ run POSTMAN
	r.POST("/resetSystem/:userID", req.ResetSystem)          //เฉพาะแอดมิน จริงๆไม่ได้เอาไปใช้ตอนโชว์ ใช้ผ่านการ run POSTMAN

	r.Run(":8090")
}
