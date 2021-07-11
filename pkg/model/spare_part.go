package model

import "time"

//SparePart ...
type SparePart struct {
	Id      uint   `gorm:"PRIMARY_KEY;" json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	InStock int    `json:"inStock"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
