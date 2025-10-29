package orm_sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// NewQueryBuilder ...
func NewQueryBuilder[T any](db *sqlx.DB, table string) QueryBuilder[T] {
	return QueryBuilder[T]{
		db:    db,
		table: table,
		Find: Read[T]{
			table: table,
			db:    db,
		},
	}
}

// Exec executes the query and returns []T
func (qb *Read[T]) Exec(ctx context.Context) ([]T, error) {
	var args []interface{}
	var whereClauses []string

	for i, cond := range qb.conditions {
		var placeholder string
		if qb.db.DriverName() == "postgres" {
			placeholder = fmt.Sprintf("$%d", i+1)
		} else { // mysql/mariadb
			placeholder = "?"
		}

		whereClauses = append(whereClauses, fmt.Sprintf("%s %s %s", cond.Field, cond.Op, placeholder))
		args = append(args, cond.Val)
	}

	query := fmt.Sprintf("SELECT * FROM %s", qb.table)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	if qb.orderBy != "" {
		query += " ORDER BY " + qb.orderBy
	}
	if qb.limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", qb.limit)
	}
	if qb.offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", qb.offset)
	}

	// Imprimir el query
	println("Executing query:", query)

	var results []T
	err := qb.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		qb.conditions = nil
		println("Error executing query:", err.Error())
		return nil, err
	}

	qb.conditions = nil // Clear conditions after execution

	return results, nil
}

// buildWhere builds the WHERE clause and args
func (qb *QueryBuilder[T]) buildWhere() (string, []interface{}) {
	var parts []string
	var args []interface{}
	for _, cond := range qb.conditions {
		parts = append(parts, fmt.Sprintf("%s %s %d", cond.Field, cond.Op, cond.Val))
		args = append(args, cond.Val)
	}
	if len(parts) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(parts, " AND "), args
}

func (qb *QueryBuilder[T]) Query(query string, params ...any) ([]T, error) {
	var results []T
	err := qb.db.Select(&results, query, params...)
	if err != nil {
		return nil, err
	}
	return results, nil
}
