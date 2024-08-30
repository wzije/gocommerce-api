package entity

type UserTable interface {
	TableName() string
}

func (User) TableName() string {
	return "users"
}

type User struct {
	BaseEntity
	Username string `gorm:"type:varchar(100);unique;not null" json:"username" validate:"required"`
	Email    string `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`
	Password string `gorm:"type:varchar(100);not null" json:"password" validate:"required"`
	Role     string `gorm:"type:varchar(20);default:'CUSTOMER'" json:"role" validate:"required,oneof=OWNER ADMIN CUSTOMER"`
	Profile  *Profile
}
