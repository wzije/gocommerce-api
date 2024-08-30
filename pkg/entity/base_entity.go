package entity

import (
	"time"
)

type BaseEntity struct {
	ID        uint64    `gorm:"primary_id;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
