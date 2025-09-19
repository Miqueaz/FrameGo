package base_controller

import (
	helpers "github.com/miqueaz/FrameGo/pkg/base/helpers"
)

// Obtener un modelo usando type assertion
func GetController[T any](name string) (*Controller[T], bool) {

	if value, ok := helpers.LoadStructure[Controller[T]](&controllers); ok {
		return value, true
	}

	return nil, false
}
