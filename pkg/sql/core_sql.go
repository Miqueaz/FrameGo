package orm_sql

import (
	"github.com/jmoiron/sqlx"
)

type Connection struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type condition struct {
	Field string
	Op    string
	Val   interface{}
}

type QueryBuilder[T any] struct {
	table      string
	db         *sqlx.DB
	conditions []condition
	orderBy    string
	limit      int
	offset     int
	Find       Read[T]
}

type Read[T any] struct {
	table      string
	db         *sqlx.DB
	conditions []condition
	orderBy    string
	limit      int
	offset     int
}
