package orm_sql

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// InsertMany inserts multiple records of type T into the table in a transaction
func (qb *QueryBuilder[T]) InsertMany(ctx context.Context, entities []T) error {
	tx, err := qb.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, entity := range entities {
		if _, err := qb.Insert(ctx, &entity); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Función principal que inserta una entidad en la base de datos
func (qb *QueryBuilder[T]) Insert(ctx context.Context, entity *T) (*T, error) {
	v := reflect.ValueOf(entity).Elem()
	t := v.Type()
	var id int

	cols, placeholders, data := buildInsertParts(t, v)
	query := buildInsertQuery(qb.table, cols, placeholders, qb.db.DriverName())
	if qb.db.DriverName() == "postgres" {
		result, err := qb.db.NamedQueryContext(ctx, query, data)
		if err != nil {
			return entity, fmt.Errorf("error executing insert query: %w", err)
		}
		defer result.Close()

		id, err = scanReturnedID(result)
	} else {
		res, err := qb.db.NamedExecContext(ctx, query, data)
		if err != nil {
			return entity, fmt.Errorf("error executing insert query: %w", err)
		}
		id64, err := res.LastInsertId()
		if err != nil {
			return entity, fmt.Errorf("error fetching last insert id: %w", err)
		}
		id = int(id64)
	}

	setEntityID(v, id)
	setEntityIDByTag(v, id, "type", "pk")
	fmt.Printf("Inserted entity with ID: %d\n", id)
	fmt.Printf("Entity after insert: %+v\n", entity)
	return entity, nil
}

// Extrae columnas, placeholders y datos desde los campos del struct
func buildInsertParts(t reflect.Type, v reflect.Value) ([]string, []string, map[string]interface{}) {
	var cols []string
	var placeholders []string
	data := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		tagSanitizer := field.Tag.Get("type")
		if tag == "" || tag == "-" || tag == "id" || tagSanitizer == "pk" {
			continue
		}
		cols = append(cols, tag)
		placeholders = append(placeholders, ":"+tag)
		data[tag] = v.Field(i).Interface()
	}

	return cols, placeholders, data
}

// Construye la sentencia SQL INSERT
func buildInsertQuery(table string, cols, placeholders []string, driver string) string {
	if driver == "postgres" {
		return fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
			table,
			strings.Join(cols, ", "),
			strings.Join(placeholders, ", "),
		)
	}

	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(cols, ", "),
		strings.Join(placeholders, ", "),
	)
}

// Escanea el ID retornado por la base de datos
func scanReturnedID(result *sqlx.Rows) (int, error) {
	var id int
	if result.Next() {
		if err := result.Scan(&id); err != nil {
			return 0, fmt.Errorf("error scanning returned id: %w", err)
		}
	}
	return id, nil
}

// Asigna el ID escaneado a la entidad (si el campo se llama ID, id o Id)
func setEntityID(v reflect.Value, id int) {
	for _, name := range []string{"ID", "id", "Id", "pk", "Pk", "PK"} {
		field := v.FieldByName(name)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.Int:
			field.SetInt(int64(id))
			return
		case reflect.Ptr:
			if field.Type().Elem().Kind() == reflect.Int {
				newID := reflect.New(field.Type().Elem())
				newID.Elem().SetInt(int64(id))
				field.Set(newID)
				return
			}
		}
	}
}

func setEntityIDByTag(v reflect.Value, id int, tagName string, tagValue string) bool {
	// 1. Nos aseguramos de que 'v' sea un struct.
	log.Printf("Setting entity ID by tag '%s:%s' to %d", tagName, tagValue, id)
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
