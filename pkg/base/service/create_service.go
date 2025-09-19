package base_service

import (
	"reflect"

	helpers "github.com/miqueaz/FrameGo/pkg/base/helpers"
	models "github.com/miqueaz/FrameGo/pkg/base/models"
)

// Crear un nuevo controlador con métodos personalizados
func NewService[T any, M any](model models.Model[M]) T {
	// Crear el valor de T usando reflect
	val := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Elem()
	// Buscar el campo 'Service' dentro de T
	serviceField := val.FieldByName("Service")
	if serviceField.IsValid() && serviceField.CanSet() {
		// Si existe el campo 'Service' y es válido, lo inicializamos
		serviceField.Set(reflect.ValueOf(Service[M]{Model: model}))
	}

	// Devolver el valor de T
	return val.Interface().(T)
}

// Guardar un controlador en el mapa global
func SaveService[T any](service *Service[T]) {
	helpers.SaveStructure(service, &services)
}
