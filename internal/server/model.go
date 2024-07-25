package server

import (
	"time"

	"github.com/gofrs/uuid"
)

type Server struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy string     `gorm:"not null;size:255" json:"created_by"`
	UpdatedBy string     `gorm:"not null;size:255" json:"updated_by"`
	DeletedBy string     `gorm:"size:255" json:"deleted_by,omitempty"`
	Name      string     `gorm:"not null;size:255" json:"name"`
	IPv4      string     `gorm:"not null;size:15" json:"ipv4"` // Assuming IPv4 standard notation
}
