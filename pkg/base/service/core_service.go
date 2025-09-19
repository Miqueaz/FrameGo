package base_service

import (
	"sync"

	"github.com/miqueaz/FrameGo/pkg/base/hooks"
	models "github.com/miqueaz/FrameGo/pkg/base/models"
)

// Controller agrupa un modelo con sus métodos
type Service[T any] struct {
	models.Model[T]
	Methods[T]
	hooks.Hookable
	hooks.Cleaners
}

type Default[T any] struct {
	Service[T]
}

// Methods define los métodos CRUD
type Methods[T any] interface {
	Read(filter map[string]any, config map[string]int) ([]T, error)
	Insert(data T) (T, error)
	Update(filter map[string]any, data T) error
	Delete(filter map[string]any) error
}

// Mapa global de modelos (uso de sync.Map para concurrencia y tipos mixtos)
var services sync.Map // key: string, value: *Model[any]
