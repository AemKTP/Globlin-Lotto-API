package models

import "encoding/json"

func UnmarshalSetting(data []byte) (Setting, error) {
	var r Setting
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Setting) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Setting struct {
	SettingID    int64 `json:"settingID"`
	WinnerPrize1 int64 `json:"winnerPrize1"`
	WinnerPrize2 int64 `json:"winnerPrize2"`
	WinnerPrize3 int64 `json:"winnerPrize3"`
	WinnerPrize4 int64 `json:"winnerPrize4"`
	WinnerPrize5 int64 `json:"winnerPrize5"`
	LotteryPrice int64 `json:"lotteryPrice"`
}
