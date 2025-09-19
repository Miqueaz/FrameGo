package connections

import "sync"

var connections sync.Map

func NewConnection[T any](name string, conn T) (bool, error) {
	_, ok := connections.Load(name)
	if ok {
		return false, nil // Connection already exists
	}

	connections.Store(name, conn)
	return true, nil // New connection created successfully
}

func GetConnection[T any](name string) (T, bool) {
	conn, ok := connections.Load(name)
	if !ok {
		var zero T
		return zero, false
	}
	c, ok := conn.(T)
	return c, ok
}
