package user

import (
	"time"
)

type User struct {
	ID        uint       `gorm:"primary_key;autoIncrement" json:"id"`
	CreatedAt time.Time  `gorm:"not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:current_timestamp on update current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy string     `gorm:"size:100;not null" json:"created_by"`
	UpdatedBy string     `gorm:"size:100;not null" json:"updated_by"`
	DeletedBy string     `gorm:"size:100" json:"deleted_by,omitempty"`
	Username  string     `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Password  string     `gorm:"not null;size:255" json:"-"`
}
