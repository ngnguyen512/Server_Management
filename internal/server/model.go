package server

import (
	"time"
)

type Server struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	CreatedBy string
	UpdatedBy string
	DeletedBy string
	Name      string
	IPv4      string
}
