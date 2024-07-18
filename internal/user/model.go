package user

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	CreatedBy string
	UpdatedBy string
	DeletedBy string
	// Add specific fields if needed, for example:
	Username string
	Password string
}
