package db

import (
	"database/sql"
)

const dbDriver = "postgres"

func NewConnection(datasource string) (*sql.DB, error) {
	conn, err := sql.Open(dbDriver, datasource)
	if err != nil {
		return nil, err
	}
	return conn, err
}
