package entity

type ProfileTable interface {
	TableName() string
}

func (Profile) TableName() string {
	return "profiles"
}

type Profile struct {
	BaseEntity
	UserID     uint64 `gorm:"not null" json:"user_id" validate:"required"`
	Name       string `gorm:"type:varchar(255);not null" json:"name" validate:"required"`
	Phone      string `gorm:"type:varchar(50)" json:"phone"`
	Address    string `gorm:"type:text" json:"address"`
	City       string `gorm:"type:varchar(255)" json:"city"`
	State      string `gorm:"type:varchar(255)" json:"state"`
	PostalCode string `gorm:"type:varchar(20)" json:"postal_code"`
	Country    string `gorm:"type:varchar(255)" json:"country"`
}
