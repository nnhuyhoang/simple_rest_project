package model

import "time"

//IssueSparePart ...
type IssueSparePart struct {
	Id          uint   `gorm:"PRIMARY_KEY;" json:"id"`
	IssueId     int    `json:"issueId"`
	SparePartId int    `json:"sparePartId"`
	Quantity    int    `json:"quantity"`
	Status      string `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	SparePart *SparePart `gorm:"foreignKey:SparePartId;" json:"sparePart"`
}
