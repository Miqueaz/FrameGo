package orm_sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Delete deletes records matching current conditions
func (qb *QueryBuilder[T]) Delete(ctx context.Context) (sql.Result, error) {
	query := fmt.Sprintf("DELETE FROM %s", qb.table)
	whereClause, _ := qb.buildWhere()
	query += whereClause
	log.Printf("Executing delete query: %s", query)
	return qb.db.ExecContext(ctx, query)
}

// DeleteByID deletes a record by its `id` field
func (qb *QueryBuilder[T]) DeleteByID(ctx context.Context, id int) (sql.Result, error) {
	qb.conditions = []condition{{Field: "id", Op: "=", Val: id}}
	return qb.Delete(ctx)
}
