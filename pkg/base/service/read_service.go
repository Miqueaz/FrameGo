package base_service

import (
	helpers "github.com/miqueaz/FrameGo/pkg/base/helpers"
)

// Obtener un modelo usando type assertion
func GetService[T any](name string) (*Service[T], bool) {

	if value, ok := helpers.LoadStructure[Service[T]](&services); ok {
		return value, true
	}

	return nil, false
}
