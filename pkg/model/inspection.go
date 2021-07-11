package model

import "time"

//Inspection ...
type Inspection struct {
	Id        uint      `gorm:"PRIMARY_KEY;" json:"id"`
	SiteId    int       `json:"siteId"`
	UserId    int       `json:"userId"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	User   *User   `gorm:"foreignKey:UserId;" json:"user"`
	Site   *Site   `gorm:"foreignKey:SiteId;" json:"site"`
	Issues []Issue `gorm:"foreignKey:InspectionId;references:Id" json:"issues"`
}
