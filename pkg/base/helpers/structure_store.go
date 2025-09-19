package base_helpers

import (
	"reflect"
	"sync"
)

func SaveStructure[T any](element *T, array *sync.Map) {
	key := reflect.TypeOf((*T)(nil)).Elem()
	array.Store(key, element)
}

func LoadStructure[T any](array *sync.Map) (*T, bool) {
	key := reflect.TypeOf((*T)(nil)).Elem()
	if value, ok := array.Load(key); ok {
		if typed, ok := value.(*T); ok {
			return typed, true
		}
	}
	return nil, false
}
