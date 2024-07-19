package postgresha

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Repository[T any] struct {
	client *ClientWrapper
	table  string
}

func NewRepository[T any](client *ClientWrapper, table string) *Repository[T] {
	return &Repository[T]{client: client, table: table}
}

func (r *Repository[T]) FindOneById(id interface{}) (T, error) {
	var result T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.table)
	err := r.client.Db().QueryRowContext(context.Background(), query, id).Scan(&result)
	return result, err
}

func (r *Repository[T]) UpdateOneById(id interface{}, items map[string]interface{}) (T, error) {
	var result T
	// Generate SQL for update statement
	setParts := []string{}
	values := []interface{}{}
	i := 1
	for key, value := range items {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", key, i))
		values = append(values, value)
		i++
	}
	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d RETURNING *", r.table, strings.Join(setParts, ", "), i)

	// Execute the update statement
	err := r.client.Db().QueryRowContext(context.Background(), query, values...).Scan(structToSlice(&result)...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository[T]) DeleteOneById(id interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.table)
	_, err := r.client.Db().ExecContext(context.Background(), query, id)
	return err
}

func structToSlice(ptr interface{}) []interface{} {
	val := reflect.ValueOf(ptr).Elem()
	results := make([]interface{}, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		results[i] = val.Field(i).Addr().Interface()
	}
	return results
}
