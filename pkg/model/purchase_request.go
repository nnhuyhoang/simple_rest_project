package model

import "time"

//PurchaseRequest ...
type PurchaseRequest struct {
	Id          uint      `gorm:"PRIMARY_KEY;" json:"id"`
	UserId      int       `json:"userId"`
	OrderDate   time.Time `json:"orderDate"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	PurchaseSpareParts []PurchaseSparePart `json:"purchaseSparePart" gorm:"foreignKey:PurchaseRequestId;references:Id"`
}
