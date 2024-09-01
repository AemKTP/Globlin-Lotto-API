package models

type Winner struct {
	WinnerID  int `gorm:"primaryKey;autoIncrement" json:"winnerID"`
	LotteryID int `gorm:"index;not null" json:"lotteryID"`
}

// Optional: Define methods if needed, such as for initialization or validation
