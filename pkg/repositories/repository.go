package repositories

import (
	"github.com/google/uuid"
)

type Repository[T any] interface {
	CreateOne(item T) (T, error)
	FindOneById(id uuid.UUID) (T, error)
	UpdateOneById(id uuid.UUID, items map[string]interface{}) (T, error)
	DeleteOneById(id uuid.UUID) error
	FindOneByAttribute(key string, value interface{}) (T, error)
}
