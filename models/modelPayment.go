package models

import (
	"encoding/json"
	"time"
)

func UnmarshalGetPayment(data []byte) (GetPayment, error) {
	var r GetPayment
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetPayment) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetPayment struct {
	PaymentID       int64     `json:"paymentID"`
	UserID          int64     `json:"userID"`
	LotteryID       int64     `json:"lotteryID"`
	TransactionType int64     `json:"transactionType"`
	Timestamp       time.Time `json:"timestamp"`
}
