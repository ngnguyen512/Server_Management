package postgresha

type Repository[T any] struct {
	client *ClientWrapper
	table  string
}

// NewRepository creates a new repository instance for a specific table.
func NewRepository[T any](client *ClientWrapper, table string) *Repository[T] {
	return &Repository[T]{client: client, table: table}
}

// FindOneById retrieves a single record by its ID.
func (r *Repository[T]) FindOneById(id interface{}) (T, error) {
	var result T
	err := r.client.Db().Table(r.table).Where("id = ?", id).First(&result).Error
	return result, err
}

// UpdateOneById updates a single record by ID and returns the updated record.
func (r *Repository[T]) UpdateOneById(id interface{}, items map[string]interface{}) (T, error) {
	var result T
	err := r.client.Db().Table(r.table).Where("id = ?", id).Updates(items).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

// DeleteOneById deletes a record by its ID.
func (r *Repository[T]) DeleteOneById(id interface{}, instance *T) error {
	result := r.client.Db().Table(r.table).Where("id = ?", id).Delete(instance)
	return result.Error
}
