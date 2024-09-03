package models

import (
	"encoding/json"

	"github.com/dgrijalva/jwt-go"
)

func UnmarshalGetCustomer(data []byte) (GetCustomer, error) {
	var r GetCustomer
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetCustomer) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetCustomer struct {
	UserID      int64  `json:"userID"`
	UserName    string `json:"userName"`
	Pwd         string `json:"Pwd"`
	Pwdconfirm  string `json:"Pwdconfirm"`
	UserType    int    `json:"userType"`
	UserBalance int64  `json:"userBalance"`
}

type GetDataCustomer struct {
	UserID      int64  `json:"userID"`
	UserName    string `json:"userName"`
	UserBalance int64  `json:"userBalance"`
	jwt.StandardClaims
}
type GetLoginCustomer struct {
	UserID   int64  `json:"userID"`
	UserName string `json:"userName"`
	UserType int    `json:"userType"`
	jwt.StandardClaims
}
