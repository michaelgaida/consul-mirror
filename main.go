package main

import (
	"os"

	"github.com/michaelgaida/consul-mirror/configuration"
	"github.com/michaelgaida/consul-mirror/consul"
	"github.com/michaelgaida/consul-mirror/storage"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Create a new instance of the logger. You can have any number of instances.
// var log = logrus.New()

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	app := cli.NewApp()
	app.Name = "consul-mirror"
	app.Author = "Michael Gaida"
	app.Copyright = "Apache License - Version 2.0, January 2004 - http://www.apache.org/licenses/"
	app.Email = "michael.gaida@protonmail.com"
	app.Description = "Mirror your consul cluster for fallback in case outages or to copy it into another environment"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "loglevel",
			Usage: "Sets the log level [debug, info, warning, error, fatal, panic]",
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
		{
			Name:  "validate",
			Usage: "Performs a basic sanity test on consul-mirror configuration files",
			UsageText: `consul-mirror validate [options] FILE
			
	Performs a basic sanity test on consul-mirror configuration files. 
	The validate command will attempt to parse the contents just as the 
	"consul-mirror" command would, and catch any errors. This is useful 
	to do a test of the configuration only, without actually starting 
	consul-mirror.

	Returns 0 if the configuration is valid, or 1 if there are problems.`,
			Action: func(c *cli.Context) {
				if c.Args().Present() {
					if err := commandValidate(c.Args().First()); err != nil {
						log.Error("Validation not succesful: ", err.Error())
						os.Exit(1)
					}
					log.Info("Validation succesful, you can go ahead and use your config")
					os.Exit(0)
				}
				cli.ShowCommandHelp(c, "validate")
			},
		},
		{
			Name:  "import",
			Usage: "import from consul",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ignoredc",
					Usage: "ignore the dc information in DB",
				},
				cli.BoolFlag{
					Name:  "dcprefix",
					Usage: "prepend dc to key, if given with setprefix this one will be prepend first -> [setprefix + dcprefix + key]",
				},
				cli.StringFlag{
					Name:  "prefix",
					Usage: "only import keys with given prefix",
				},
				cli.StringFlag{
					Name:  "setprefix",
					Usage: "prepend given prefix to key, if given with dcprefix this one will be prepend last -> [setprefix + dcprefix + key]",
				},
			},
			Action: func(c *cli.Context) {
				storage, consul, err := initConsul(c)
				defer storage.Close()
				if err != nil {
					log.Error("Failed importing data: ", err.Error())
				}
				commandImport(storage, consul, c.BoolT("ignoredc"), c.BoolT("dcprefix"), c.String("prefix"), c.String("setprefix"))
			},
		},
		{
			Name:  "export",
			Usage: "export from consul",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "ignoredc",
					Usage: "ignore the original dc",
				},
				cli.StringFlag{
					Name:  "dc",
					Usage: "only export from given dc",
				},
				cli.BoolFlag{
					Name:  "incremental",
					Usage: "Only create a new entry when the value (for the DC) changed",
				},
				cli.StringFlag{
					Name:  "prefix",
					Usage: "key prefix for keys to be exported",
				},
			},
			Action: func(c *cli.Context) {
				storage, consul, err := initConsul(c)
				defer storage.Close()
				if err != nil {
					log.Error("Failed exporting data: ", err.Error())
				}
				commandExport(storage, consul, c.BoolT("ignoredc"), c.BoolT("incremental"), c.String("dc"), c.String("prefix"))
			},
		},
	}

	app.Run(os.Args)
}

func initConsul(cli *cli.Context) (*storage.Mssql, *consul.Consul, error) {
	config := configuration.GetConfig("config.json")

	config.OverwriteConfig(cli)
	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Panic("Cant parse log level: ", err.Error())
		return nil, nil, err
	}
	log.SetLevel(logLevel)
	log.Debug(config.PrintDebug())

	// s := storage.Mssql{}
	conn, err := storage.OpenConnection(config)
	if err != nil {
		log.Error("Error initializing storage")
		return nil, nil, err
	}

	consul := consul.GetConsul(config)

	return conn, consul, nil
}

func commandExport(storage *storage.Mssql, consul *consul.Consul, ignoreDC, incremental bool, dc, prefix string) {
	dcs, err := consul.GetDCs(dc)
	if err != nil {
		log.Error("commandExport: Error getting the DCs from Consul: ", err.Error())
		os.Exit(1)
	}
	kvs, err := consul.GetKVs(prefix, dcs)
	if err != nil {
		log.Error("commandExport: Error getting the KVs from Consul: ", err.Error())
		os.Exit(1)
	}
	err = storage.WriteKVs(kvs, ignoreDC, incremental)
	if err != nil {
		log.Error("commandExport: Error writing the KVs to storage: ", err.Error())
		os.Exit(1)
	}
}

func commandImport(storage *storage.Mssql, consul *consul.Consul, ignoreDC, dcprefix bool, prefix, setprefix string) {
	// log.Fatal("IMPORT!!!!!!!")
	kvs, err := storage.GetKVs(prefix)
	if err != nil {
		log.Error("commandImport: Error getting data from the storage: ", err.Error())
		os.Exit(1)
	}
	err = consul.WriteKVs(kvs, ignoreDC, dcprefix, setprefix)
	if err != nil {
		log.Error("commandImport: Error writing KVs to consul: ", err.Error())
		os.Exit(1)
	}
}

func commandValidate(file string) error {
	testConfiguration := configuration.GetConfig(file)

	err := testConfiguration.ValidateConfiguration()
	if err != nil {
		return err
	}
	return nil
}
