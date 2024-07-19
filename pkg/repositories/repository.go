package repositories

type Repository[T any] interface {
	CreateOneById(id interface{}, items map[string]interface{}) (T, error)
	FindOneById(id interface{}) (T, error)
	UpdateOneById(id interface{}, items map[string]interface{}) (T, error)
	DeleteOneById(id interface{}) error
}
