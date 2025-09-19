package base_models

import (
	"log"
	helpers "main/pkg/base/helpers"
)

// Obtener un modelo usando type assertion
func GetModel[T any]() (*Model[T], bool) {
	log.Print("Cargando Modelo")

	if value, ok := helpers.LoadStructure[Model[T]](&models); ok {
		log.Print("Modelo Obtenido: ", value.Name)
		return value, true
	}
	log.Print("No se encontro el modelo")
	return nil, false
}
