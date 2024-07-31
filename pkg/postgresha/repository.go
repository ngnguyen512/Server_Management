package postgresha

import (
	"github.com/google/uuid"
)

type Repository[T any] struct {
	client *ClientWrapper
}

// NewRepository creates a new repository instance.
func NewRepository[T any](client *ClientWrapper) *Repository[T] {
	return &Repository[T]{client: client}
}

// FindOneById retrieves a single record by its ID.
func (r *Repository[T]) FindOneById(id uuid.UUID) (T, error) {
	var result T
	err := r.client.Db().Model(&result).Where("id = ?", id).First(&result).Error
	return result, err
}

// UpdateOneById updates a single record by ID and returns the updated record.
func (r *Repository[T]) UpdateOneById(id uuid.UUID, items map[string]interface{}) (T, error) {
	var result T
	err := r.client.Db().Model(&result).Where("id = ?", id).Updates(items).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

// DeleteOneById deletes a record by its ID.
func (r *Repository[T]) DeleteOneById(id uuid.UUID) error {
	var instance T
	result := r.client.Db().Model(&instance).Where("id = ?", id).Delete(&instance)
	return result.Error
}

// CreateOne inserts a new item into the database.
func (r *Repository[T]) CreateOne(item T) (T, error) {
	err := r.client.Db().Model(&item).Create(&item).Error
	if err != nil {
		return item, err
	}
	return item, nil
}
func (r *Repository[T]) FindOneByAttribute(field string, value interface{}) (T, error) {
	var result T
	err := r.client.Db().Model(&result).Where(field+" = ?", value).First(&result).Error
	return result, err
}
