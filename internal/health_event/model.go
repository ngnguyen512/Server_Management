package health_event

import (
	"time"

	"github.com/google/uuid"
)

type HealthEvent struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time  `json: "created_at" gorm:"column:created_at; not null;"`
	UpdatedAt time.Time  `json: "updated_at" gorm:"column:updated_at; not null;"`
	DeletedAt *time.Time `json: "deleted_at" gorm:"column:deleted_at; not null;"`
	CreatedBy string     `json: "created_by" gorm:"column:created_at; not null;"`
	UpdatedBy string     `json: "updated_by"`
	DeletedBy string     `json: "deleted_by"`
	Server_id int        `json: "server_id"`
	Status    int        `json: "status"`
	// Add specific fields if needed
}
