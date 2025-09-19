package orm_sql

import (
	"database/sql"
	"fmt"
)

func InitPostgres(connection Connection) (*sql.DB, error) {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		connection.Host,
		connection.Port,
		connection.User,
		connection.Password,
		connection.Database,
		connection.SSLMode,
	)

	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = DB.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("[Postgres] Connection Successful")

	return DB, nil
}
