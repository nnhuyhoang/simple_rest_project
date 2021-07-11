package model

import (
	"time"

	"gorm.io/datatypes"
)

// UserActionLog --
type UserActionLog struct {
	Id         uint           `gorm:"PRIMARY_KEY;" json:"id"`
	UserID     int            `json:"userId"`
	TargetId   int            `json:"targetId"`
	TargetType string         `json:"targetType"`
	Action     string         `json:"action"`
	OldData    datatypes.JSON `sql:"type:JSONB DEFAULT '{}'::JSONB" json:"oldData"`
	NewData    datatypes.JSON `sql:"type:JSONB DEFAULT '{}'::JSONB" json:"newData"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// EmptyObject --
type EmptyObject struct{}
