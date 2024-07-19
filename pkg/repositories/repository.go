package repositories

type Repository[T any] interface {
	FindOneById(id interface{}) (T, error)
	UpdateOneById(id interface{}, items map[string]interface{}) (T, error)
	DeleteOneById(id interface{}) error
}
