package storage

// Storage defines the interface that should be implemented
// if it is needed to add a new storage type.
type Storage interface {
	Store(data *map[string]interface{})
	Get(key string) (interface{}, bool)
	GetAll() map[string]interface{}
}
