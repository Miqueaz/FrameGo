package orm_sql

// Where adds a condition
func (qb *Read[T]) Where(field, op string, val interface{}) *Read[T] {
	qb.conditions = append(qb.conditions, condition{field, op, val})
	return qb
}

// Limit sets a limit
func (qb *Read[T]) Limit(n int) *Read[T] {
	qb.limit = n
	return qb
}

// Offset sets an offset
func (qb *Read[T]) Offset(n int) *Read[T] {
	qb.offset = n
	return qb
}

// OrderBy sets order
func (qb *Read[T]) OrderBy(clause string) *Read[T] {
	qb.orderBy = clause
	return qb
}
