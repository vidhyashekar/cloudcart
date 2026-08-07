package model

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (RefreshToken) TableName() string {
	return "auth.refresh_tokens"
}
