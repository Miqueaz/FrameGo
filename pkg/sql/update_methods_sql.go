// update.go
package orm_sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Update updates fields of a record matching current conditions based on non-zero values of entity
func (qb *QueryBuilder[T]) Update(ctx context.Context, entity T) (sql.Result, error) {
	v := reflect.ValueOf(entity)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	var sets []string
	data := map[string]interface{}{}
	paramIdx := 1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" || tag == "-" || tag == "id" {
			continue
		}
		val := v.Field(i).Interface()
		// skip zero values
		if reflect.DeepEqual(val, reflect.Zero(field.Type).Interface()) {
			continue
		}
		sets = append(sets, fmt.Sprintf("%s = :%s", tag, tag))
		data[tag] = val
		paramIdx++
	}
	if len(sets) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	query := fmt.Sprintf("UPDATE %s SET %s", qb.table, strings.Join(sets, ", "))
	// add where
	whereClause, _ := qb.buildWhere()
	query += whereClause

	log.Printf("Where clause: %s\n", whereClause)
	return qb.db.NamedExecContext(ctx, query, data)
}

// UpdateByID updates a record by its `id` field based on non-zero values of entity
func (qb *QueryBuilder[T]) UpdateByID(ctx context.Context, id int, entity T, field string) (sql.Result, error) {
	qb.conditions = []condition{{Field: field, Op: "=", Val: id}}
	log.Printf("Updating record with ID %d in table %s with data: %+v\n", id, qb.table, entity)
	result, err := qb.Update(ctx, entity)
	return result, err
}
