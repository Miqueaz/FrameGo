package base_controller

import (
	helpers "github.com/miqueaz/GoRestAPI/pkg/base/helpers"
	base_service "github.com/miqueaz/GoRestAPI/pkg/base/service"
)

// Crear un nuevo controlador
func NewController[T any](s base_service.Service[T], m Methods[T]) *Controller[T] {
	if m == nil {
		controller := &Controller[T]{Service: s, Methods: &Controller[T]{}}
		SaveController(controller)
		return controller
	}

	controller := &Controller[T]{Service: s, Methods: m}
	SaveController(controller)
	return controller
}

// Guardar un controlador en el mapa global
func SaveController[T any](controller *Controller[T]) {
	helpers.SaveStructure(controller, &controllers)
}
