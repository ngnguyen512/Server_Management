package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time      `json: "created_at"`
	UpdatedAt time.Time      `json: "updated_at"`
	DeletedAt gorm.DeletedAt `json: "deleted_at"`
	CreatedBy string         `json: "created_by"`
	UpdatedBy string         `json: "updated_by"`
	DeletedBy string         `json: "deleted_by"`
	Username  string         `json: "username" validate:"required"`
	Password  string         `json: "password" validate:"required"`
}
