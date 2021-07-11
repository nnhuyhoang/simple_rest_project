package model

import "time"

//Role ...
type Role struct {
	Id   uint   `gorm:"PRIMARY_KEY;" json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
