package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/michaelgaida/consul-mirror/configuration"
	"github.com/michaelgaida/consul-mirror/consul"
	"github.com/michaelgaida/consul-mirror/storage"
)

func main() {
	var (
		validate = flag.String("validate", "", "configuration file to be validated")
	)

	flag.Parse()

	if *validate != "" {
		// check the new config
		testConfiguration := configuration.GetConfig(*validate)
		os.Exit(testConfiguration.ValidateConfiguration())
	}

	config := configuration.GetConfig("config.json")
	if config.Debug {
		fmt.Printf(config.PrintDebug())
	}

	// s := storage.Mssql{}
	conn := storage.OpenConnection(config)
	defer conn.Close()

	consul := consul.GetConsul(config)
	kvs := consul.GetKVs("")
	if kvs == nil {
	}
	conn.WriteKVs(kvs)
	services := consul.GetServices()
	if services == nil {

	}
	nodes := consul.GetNodes()
	if nodes == nil {

	}
}
