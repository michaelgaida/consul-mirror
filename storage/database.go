package storage

import (
	"database/sql"

	"github.com/michaelgaida/consul-mirror/configuration"
)

type storage interface {
	OpenConnection(configuration configuration.Struct) *sql.DB
}

type Mssql struct {
}

// OpenConnection opens a database connection and returns it
func (*Mssql) OpenConnection(config *configuration.Struct) *sql.DB {
	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", config.Database.Server, config.Database.User, config.Database.Password, config.Database.Port, config.Database.Database)
	// // if config.Debug {
	// // 	fmt.Printf(" connString:%s\n", connString)
	// // }

	// conn, err := sql.Open("mssql", connString)
	// if err != nil {
	// 	log.Fatal("Open connection failed:", err.Error())
	// }
	// return conn
	return nil
}
