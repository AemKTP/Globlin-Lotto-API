package models

import "encoding/json"

func UnmarshalGetLottery(data []byte) (GetLottery, error) {
	var r GetLottery
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetLottery) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetLottery struct {
	LotteryID     int64  `json:"lotteryID"`
	LotteryNumber string `json:"lotteryNumber"`
}
