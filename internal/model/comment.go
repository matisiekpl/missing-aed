package model

import "time"

type Comment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	NodeID    string    `gorm:"index;not null" json:"nodeID"`
	Content   string    `gorm:"type:varchar(200);not null" json:"content"`
}
