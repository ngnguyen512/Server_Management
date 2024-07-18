package health_event

import (
	"time"
)

type HealthEvent struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	CreatedBy string
	UpdatedBy string
	DeletedBy string
	Server_id int
	Status    int
	// Add specific fields if needed
}
