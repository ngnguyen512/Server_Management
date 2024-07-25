package health_event

import (
	"time"
)

type HealthEvent struct {
	ID        uint       `gorm:"primary_key;autoIncrement"`                                      // Ensures ID is the primary key and automatically increments
	CreatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`                             // Ensures field is not nullable and has a default value
	UpdatedAt time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // Updates time automatically on update
	DeletedAt *time.Time `gorm:"index;default:NULL"`                                             // Allows NULL, useful for soft delete functionality
	CreatedBy string     `gorm:"size:255;not null"`                                              // Specifies maximum size and not null
	UpdatedBy string     `gorm:"size:255;not null"`
	DeletedBy string     `gorm:"size:255;default:NULL"` // Allows NULL values
	Server_id int        `gorm:"not null;index"`        // Adds an index for quicker queries
	Status    int        `gorm:"not null"`              // Ensures status is not nullable
	// Add specific fields if needed
}
