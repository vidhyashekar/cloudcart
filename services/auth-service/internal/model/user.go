package model

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	FirstName    string `gorm:"size:100;not null"`
	LastName     string `gorm:"size:100"`
	Email        string `gorm:"size:255;unique;not null"`
	PasswordHash string `gorm:"not null"`

	RoleID uint
	Role   Role `gorm:"foreignKey:RoleID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "auth.users"
}
