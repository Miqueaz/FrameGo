package orm_sql

import (
	"context"
	"fmt"
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

	cols, placeholders, data := buildInsertParts(t, v)
	query := buildInsertQuery(qb.table, cols, placeholders)

	result, err := qb.db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return entity, fmt.Errorf("error executing insert query: %w", err)
	}
	defer result.Close()

	id, err := scanReturnedID(result)
	if err != nil {
		return entity, err
	}

	setEntityID(v, id)
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
		if tag == "" || tag == "-" || tag == "id" {
			continue
		}
		cols = append(cols, tag)
		placeholders = append(placeholders, ":"+tag)
		data[tag] = v.Field(i).Interface()
	}

	return cols, placeholders, data
}

// Construye la sentencia SQL INSERT
func buildInsertQuery(table string, cols, placeholders []string) string {
	return fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING id",
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
	for _, name := range []string{"ID", "id", "Id"} {
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
