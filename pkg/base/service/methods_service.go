package base_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	base_helpers "github.com/miqueaz/FrameGo/pkg/base/helpers"
)

// Método Read con soporte de hooks
func (s *Service[T]) Read(filter map[string]any) ([]T, error) {

	// Ejecutar hooks antes del Read
	if err := s.ExecuteHooks(s.BeforeRead, filter); err != nil {
		return nil, err
	}

	//Imprimir los filtros
	fmt.Printf("Leyendo datos de la colección '%s' con el filtro: %v\n", s.Model.CollectionName, filter)
	filtros := base_helpers.NormalizarFiltros(filter)

	qb := s.Model.Find
	// qb.OrderBy("ID DESC")
	for campo, cond := range filtros {
		for i := 0; i < len(cond); i += 2 {
			op := fmt.Sprintf("%v", cond[i])
			val := cond[i+1]
			qb = *qb.Where(campo, op, val)
		}
	}
	// Agregar sort

	data, err := qb.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	if data == nil {
		err = errors.New("no data found")
	}

	// Ejecutar hooks después del Read
	if err == nil {
		_ = s.ExecuteHooks(s.AfterRead, filter) // Ignoramos errores de hooks posteriores
	}

	return data, err
}

// Insert: Insertar datos en la base de datos
func (s *Service[T]) Insert(data T) (T, error) {
	fmt.Printf("Insertando datos en la colección '%s': %v\n", data)
	result, err := s.Model.Insert(context.Background(), &data)
	return *result, err
}

// Update: Actualizar datos en la base de datos
func (s *Service[T]) Update(idStr string, data T) (T, error) {
	id, err := strconv.Atoi(idStr)
	fmt.Printf("Actualizando datos de la colección '%s' con el filtro: %v y los datos: %v\n", id, data)
	// return db.UpdateDocument(filter, data, b.Model.CollectionName) // Implementación real de actualización
	_, err = s.Model.UpdateByID(context.Background(), id, data)
	// Agregar a la data el ID actualizado con un set
	v := reflect.ValueOf(&data).Elem()
	if idField := v.FieldByName("ID"); idField.IsValid() && idField.CanSet() {
		idField.SetInt(int64(id))
	}
	return data, err

}

// Delete: Eliminar datos de la base de datos
func (s *Service[T]) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	fmt.Printf("Eliminando datos de la colección '%s' con el filtro: %v\n", id)
	// return db.DeleteDocument(filter, b.Model.CollectionName) // Implementación real de eliminación
	_, err = s.Model.DeleteByID(context.Background(), id)
	return err

}

func (s *Service[T]) ReadOne(id int) (T, error) {
	// Implementation for reading user data
	if id == 0 {
		return *new(T), errors.New("ID cannot be zero")
	}

	data, err := s.Model.Find.Where("ID", "=", id).Exec(context.Background())
	if err != nil {
		return *new(T), err
	}

	if len(data) == 0 {
		return *new(T), errors.New("no data found")
	}

	return data[0], err
}

func Sanitizar[S any, E any](data E) (S, error) {
	// Inicializar un valor vacío de tipo S
	var sanitized S

	// Convertir los datos a JSON para sanitización (por ejemplo, eliminando valores peligrosos)
	sanitizedData, err := json.Marshal(data)
	if err != nil {
		return sanitized, errors.New("failed to sanitize data: " + err.Error())
	}

	// En este punto, podrías modificar o verificar los datos en sanitizedData
	// Aquí se puede agregar lógica para procesar los datos antes de deserializarlos.

	// Deserializar los datos sanitizados de vuelta al tipo `S`
	err = json.Unmarshal(sanitizedData, &sanitized)
	if err != nil {
		return sanitized, errors.New("failed to unmarshal sanitized data: " + err.Error())
	}

	// Retornar los datos sanitizados y sin errores
	return sanitized, nil
}
