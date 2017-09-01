package main

import (
	"log"
	"os"

	"github.com/michaelgaida/consul-mirror/configuration"
	"github.com/michaelgaida/consul-mirror/consul"
	"github.com/michaelgaida/consul-mirror/storage"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "consul-mirror"
	app.Author = "Michael Gaida"
	app.Copyright = "Apache License - Version 2.0, January 2004 - http://www.apache.org/licenses/"
	app.Email = "michael.gaida@protonmail.com"
	app.Description = "Mirror your consul cluster for fallback in case outages or to copy it into another environment"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Switch on the verbose mode",
		},
		cli.StringFlag{
			Name:  "token",
			Usage: "ACL Token to be used to interact with consul",
		},
		cli.StringFlag{
			Name:  "dbpassword",
			Usage: "Database password to be used to connect",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name: "validate",
			UsageText: `consul-mirror validate [options] FILE
			
	Performs a basic sanity test on consul-mirror configuration files. 
	The validate command will attempt to parse the contents just as the 
	"consul-mirror" command would, and catch any errors. This is useful 
	to do a test of the configuration only, without actually starting 
	consul-mirror.

	Returns 0 if the configuration is valid, or 1 if there are problems.`,
			Action: func(c *cli.Context) {
				if c.Args().Present() {
					os.Exit(commandValidate(c.Args().First()))
				}
				cli.ShowCommandHelp(c, "validate")
			},
		},
		cli.Command{
			Name:  "import",
			Usage: "import from consul",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ignoredc",
					Usage: "ignore the dc information in DB",
				},
				cli.StringFlag{
					Name:  "prefix",
					Usage: "key prefix for keys to be imported",
				},
			},
			Action: func(c *cli.Context) {
				storage, consul := initConsul(c)
				defer storage.Close()
				commandImport(storage, consul, c.BoolT("ignoredc"), c.String("prefix"))
			},
		},
		cli.Command{
			Name:  "export",
			Usage: "export from consul",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ignoredc",
					Usage: "ignore the original dc",
				},
				cli.BoolFlag{
					Name:  "incversion",
					Usage: "creates a new version for every KV",
				},
				cli.StringFlag{
					Name:  "prefix",
					Usage: "key prefix for keys to be exported",
				},
			},
			Action: func(c *cli.Context) {
				storage, consul := initConsul(c)
				defer storage.Close()
				commandExport(storage, consul, c.BoolT("ignoredc"), c.BoolT("incversion"), c.String("prefix"))
			},
		},
	}

	app.Run(os.Args)
}

func initConsul(cli *cli.Context) (*storage.Mssql, *consul.Consul) {
	config := configuration.GetConfig("config.json")

	config.OverwriteConfig(cli)

	if config.Debug {
		log.Printf(config.PrintDebug())
	}

	// s := storage.Mssql{}
	conn := storage.OpenConnection(config)

	consul := consul.GetConsul(config)

	return conn, consul
}

func commandExport(storage *storage.Mssql, consul *consul.Consul, ignoreDC, incversion bool, prefix string) {
	dcs := consul.GetDCs()
	kvs := consul.GetKVs(prefix, dcs)
	storage.WriteKVs(kvs, ignoreDC, incversion)
}

func commandImport(storage *storage.Mssql, consul *consul.Consul, ignoreDC bool, prefix string) {
	kvs, err := storage.GetKVs(prefix)
	if err != nil {
		log.Fatal("Error while fetching the data from the storage: ", err)
	}
	err = consul.WriteKVs(kvs, ignoreDC)
	if err != nil {
		log.Fatal("Error while writing the data to consul", err)
	}
}

func commandValidate(file string) int {
	testConfiguration := configuration.GetConfig(file)
	return testConfiguration.ValidateConfiguration()
}
