package model

import "time"

//User ..
type User struct {
	Id             uint      `gorm:"PRIMARY_KEY;" json:"Id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	FullName       string    `json:"fullName"`
	Email          string    `json:"email"`
	PhoneNumber    string    `json:"phoneNumber"`
	HashedPassword string    `json:"hashedPassword"`
	RoleId         int       `json:"roleId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	Role *Role `gorm:"foreignKey:RoleId;" json:"role"`
}
