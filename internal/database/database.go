package database

// Database defines an interface for our data access layer.
type Database interface {
	Health() interface{}
}

// InMemoryDB is a dummy implementation of the Database interface.
type InMemoryDB struct{}

// NewInMemoryDB creates and returns a new InMemoryDB instance.
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{}
}

// Health returns a dummy health status.
func (db *InMemoryDB) Health() interface{} {
	return map[string]string{"status": "healthy"}
}
