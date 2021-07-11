package model

import "time"

//PurchaseSparePart ...
type PurchaseSparePart struct {
	Id                uint      `gorm:"PRIMARY_KEY;" json:"id"`
	PurchaseRequestId int       `json:"purchaseRequestId"`
	SparePartId       int       `json:"sparePartId"`
	Quantity          int       `json:"quantity"`
	SparePartName     string    `json:"sparePartName"`
	SparePartCode     string    `json:"sparePartCode"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
