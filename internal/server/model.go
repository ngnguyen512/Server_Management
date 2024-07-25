package server

import (
	"time"

	"github.com/gofrs/uuid"
)

type Server struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	CreatedBy string
	UpdatedBy string
	DeletedBy string
	Name      string `json:"status"`
	IPv4      string
}
