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

	r.GET("/lotteries", res.GetLotterys)
	r.GET("/canbuylotteries", res.GetCanBuyLotteries)
	r.GET("/MyLottery/:userID", res.GetMyLottery)
	r.GET("/profile/:userID", res.GetProfile)
	r.GET("/lotteriesSearch/:lotterynumber", res.GetlotteriesSearch)
	r.GET("/CheckLotteryResult/:lotteryResult", res.GetCheckLotteryResult)
	r.GET("/AllLotteryResult/", res.GETAllLotteryResults) // lottery ที่ออกรางวัล

	// Admin
	r.GET("/users", res.GetUsers)
	r.GET("/randomlotteryResult", res.RandomResult) //เฉพาะแอดมิน จริงๆไม่ได้เอาไปใช้ตอนโชว์ ใช้ผ่านการ run POSTMAN

	r.POST("/register", req.Register)
	r.POST("/login", req.Login)
	r.POST("/buylottery/:userID", req.BuyLottery)
	r.POST("/cashin/:userID", req.CashIn)

	//เฉพาะแอดมิน จริงๆไม่ได้เอาไปใช้ตอนโชว์ ใช้ผ่านการ run POSTMAN
	r.POST("/randomlotterynumber", req.Random)

	r.Run(":8090")
}
