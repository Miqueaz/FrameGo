package base_models

import (
	ORM "main/pkg/sql"
	"sync"
)

// Modelo genérico
type Model[T any] struct {
	ID             int
	Name           string
	CollectionName string
	Structure      T
	ORM.QueryBuilder[T]
}

// Mapa global de modelos (uso de sync.Map para concurrencia y tipos mixtos)
var models sync.Map // key: string, value: *Model[any]
