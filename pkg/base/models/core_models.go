package base_models

import (
	"sync"

	ORM "github.com/miqueaz/FrameGo/pkg/sql"
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
