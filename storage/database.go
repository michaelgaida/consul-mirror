package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/michaelgaida/consul-mirror/configuration"
)

type storage interface {
	OpenConnection(configuration configuration.Struct) *sql.DB
}

type Mssql struct {
	conn  *sql.DB
	debug bool
}

// OpenConnection opens a database connection and returns it
func OpenConnection(config *configuration.Struct) *Mssql {
	var result Mssql

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Port, config.Database.Database)
	if config.Debug {
		fmt.Printf(" connString:%s\n", connString)
	}

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	result.conn = conn
	result.debug = config.Debug
	return &result
}

func (m *Mssql) Close() {
	m.Close()
}
