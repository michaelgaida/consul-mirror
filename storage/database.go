package storage

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	// This comment is for golint since it does not recognice the _
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/michaelgaida/consul-mirror/configuration"
)

type storage interface {
	OpenConnection(configuration configuration.Struct) *sql.DB
}

// Mssql is the struct representing the Mssql connection
type Mssql struct {
	conn *sql.DB
}

// OpenConnection opens a database connection and returns it
func OpenConnection(config *configuration.Struct) (*Mssql, error) {
	var result Mssql

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Port, config.Database.Database)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Error("Open connection failed: ", err.Error())
		return nil, err

	}
	result.conn = conn
	return &result, nil
}

// Close the Mssql connection
func (m *Mssql) Close() {
	m.conn.Close()
}
