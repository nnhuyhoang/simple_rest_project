package model

import "time"

//Site ...
type Site struct {
	Id          uint   `gorm:"PRIMARY_KEY;" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Status      string `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Inspections []Inspection `gorm:"foreignKey:SiteId;references:Id" json:"inspections"`
	Issues      []Issue      `gorm:"foreignKey:SiteId;references:Id" json:"issues"`
}
