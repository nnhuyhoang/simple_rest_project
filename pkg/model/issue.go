package model

import "time"

//Issue ...
type Issue struct {
	Id           uint   `gorm:"PRIMARY_KEY;" json:"id"`
	SiteId       int    `json:"siteId"`
	UserId       int    `json:"userId"`
	InspectionId int    `json:"inspectionId"`
	Status       string `json:"status"`
	Description  string `json:"description"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	IssueSpareParts []IssueSparePart `json:"issueSparePart" gorm:"foreignKey:IssueId;references:Id"`
}
