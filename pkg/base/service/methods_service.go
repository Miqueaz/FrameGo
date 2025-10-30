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
	v := reflect.ValueOf(s.Model.Structure)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return data, fmt.Errorf("s.Model no es un struct (es %s)", v.Kind())
	}

	var field string
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		if structField.Tag.Get("type") == "pk" {
			field = structField.Tag.Get("db")
			break
		}
	}
	if field == "" {
		return data, fmt.Errorf("no se encontró ningún campo con el tag 'type:\"pk\"'")
	}

	fmt.Printf("Actualizando registro con ID: %d usando el campo: %s\n", id, field)
	// return db.UpdateDocument(filter, data, b.Model.CollectionName) // Implementación real de actualización
	_, err = s.Model.UpdateByID(context.Background(), id, data, field)
	// Agregar a la data el ID actualizado con un set
	v2 := reflect.ValueOf(&data).Elem()
	setEntityIDByTag(v2, id, "type", "pk")
	return data, err

}

// Delete: Eliminar datos de la base de datos
func (s *Service[T]) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Es bueno manejar el error de conversión aquí
		return fmt.Errorf("ID inválido '%s': %w", idStr, err)
	}

	// 1. Obtenemos el reflect.Value
	v := reflect.ValueOf(s.Model.Structure)

	// 2. [LA MEJORA] Hacemos la lógica robusta
	// Si 'v' es un puntero (reflect.Ptr), usamos .Elem()
	// para obtener el struct al que apunta.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 3. Ahora nos aseguramos de que lo que tenemos es un struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("s.Model no es un struct (es %s)", v.Kind())
	}

	// 4. Tu lógica de bucle (que ya estaba correcta)
	var field string
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field %d: %s\n", i, v.Type().Field(i).Name)
		structField := v.Type().Field(i)
		if structField.Tag.Get("type") == "pk" {
			// Asumo que quieres el nombre de la columna de la DB
			field = structField.Tag.Get("db")
			break
		}
	}

	// 5. Es bueno verificar si encontramos el campo
	if field == "" {
		return fmt.Errorf("no se encontró ningún campo con el tag 'type:\"pk\"'")
	}

	// 6. Corrección del Printf:
	// 'id' es un 'int' (se usa %d), no un string (%s).
	// También añadí el 'field' que encontramos.
	fmt.Printf("Eliminando ID: %d usando el campo PK: %s\n", id, field)

	// 7. Ejecutar la eliminación
	// Esta línea asume que tu tipo T tiene este método.
	_, err = s.Model.DeleteByID(context.Background(), id, field)
	return err
}

func (s *Service[T]) ReadOne(id int) (T, error) {
	// Implementation for reading user data
	if id == 0 {
		return *new(T), errors.New("ID cannot be zero")
	}
	// Obtener el reflect.Value
	v := reflect.ValueOf(s.Model.Structure)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return *new(T), errors.New("s.Model is not a struct")
	}

	var field string
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		if structField.Tag.Get("type") == "pk" {
			field = structField.Tag.Get("db")
			break
		}
	}
	if field == "" {
		return *new(T), errors.New("no field with tag 'type:\"pk\"' found")
	}

	fmt.Printf("Reading one record with ID: %d using field: %s\n", id, field)

	// Construir el query para buscar por el campo PK

	data, err := s.Model.Find.Where(field, "=", id).Exec(context.Background())
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

func setEntityIDByTag(v reflect.Value, id int, tagName string, tagValue string) bool {
	// 1. Nos aseguramos de que 'v' sea un struct.
	fmt.Printf("Setting entity ID by tag '%s:%s' to %d", tagName, tagValue, id)
	if v.Kind() != reflect.Struct {
		return false
	}

	// 2. Iteramos sobre todos los campos del struct.
	for i := 0; i < v.NumField(); i++ {

		// 3. Obtenemos la especificación del campo (para leer el tag).
		// v.Type() nos da el "molde" del struct.
		fieldSpec := v.Type().Field(i)

		// 4. Leemos el tag específico.
		tag := fieldSpec.Tag.Get(tagName)

		// 5. ¡Comprobación! Verificamos si es el tag que buscamos.
		if tag == tagValue {

			// 6. Obtenemos el valor de ese campo (para modificarlo).
			fieldValue := v.Field(i)

			// 7. Verificamos que se pueda modificar.
			if !fieldValue.IsValid() || !fieldValue.CanSet() {
				continue // No se puede, probemos el siguiente (aunque PK suele ser única)
			}

			// 8. Usamos la misma lógica de antes para asignar el valor.
			switch fieldValue.Kind() {
			case reflect.Int:
				fieldValue.SetInt(int64(id))
				return true // ¡Éxito!
			case reflect.Ptr:
				if fieldValue.Type().Elem().Kind() == reflect.Int {
					newID := reflect.New(fieldValue.Type().Elem())
					newID.Elem().SetInt(int64(id))
					fieldValue.Set(newID)
					return true // ¡Éxito!
				}
			}

			// Si encontramos el tag "pk" pero el tipo no era int o *int
			// (ej. era un string), la asignación falla y retornamos false.
			return false
		}
	}

	// 9. Si el bucle termina, no se encontró ningún campo con ese tag.
	return false
}
