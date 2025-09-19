package base_controller

import (
	"sync"

	base_service "github.com/miqueaz/FrameGo/pkg/base/service"

	"github.com/gin-gonic/gin"
)

// Controller agrupa un modelo con sus métodos
type Controller[T any] struct {
	base_service.Service[T]
	Methods[T]
}

// Methods define los métodos CRUD
type Methods[T any] interface {
	Read(gin *gin.Context)
	Insert(gin *gin.Context)
	Update(gin *gin.Context)
	Delete(gin *gin.Context)
}

// Mapa global de modelos (uso de sync.Map para concurrencia y tipos mixtos)
var controllers sync.Map // key: string, value: *Model[any]
