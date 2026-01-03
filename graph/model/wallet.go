package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Address string `json:"address" gorm:"unique"`
	Tokens  int    `json:"tokens"`
}
