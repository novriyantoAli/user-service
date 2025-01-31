package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID `gorm:"type:uuid;not null"`
	Name      string    `gorm:"varchar(100);not null"`
	Password  string    `gorm:"varcher(255);not null"`
	Phone     string    `gorm:"varchar(15);not null"`
	Email     string    `gorm:"varcher(100);not null"`
	RoleID    uint      `gorm:"type:uint;not null"`
	CreatedAt *time.Time
	UpdateAt  *time.Time
	Role      Role `gorm:"foreignKey:role_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
