package health_event

import (
	"time"

	"github.com/google/uuid"
)

type HealthEvent struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
	UpdatedBy string
	DeletedBy string
	Server_id int
	Status    int
	// Add specific fields if needed
}
