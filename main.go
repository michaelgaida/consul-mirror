package main

import (
	"fmt"
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
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Switch on the verbose mode",
		},
	}

	validateHelpText := `consul-mirror validate [options] FILE
	
		Performs a basic sanity test on consul-mirror configuration files. 
		The validate command will attempt to parse the contents just as the 
		"consul-mirror" command would, and catch any errors. This is useful 
		to do a test of the configuration only, without actually starting 
		consul-mirror.
	
		Returns 0 if the configuration is valid, or 1 if there are problems.`
	app.Commands = []cli.Command{
		cli.Command{
			Name:      "validate",
			UsageText: validateHelpText,

			Action: func(c *cli.Context) {
				commandValidate(c.Args().First(), validateHelpText)
			},
		},
		cli.Command{
			Name:  "import, i",
			Usage: "import from consul",
			Action: func(c *cli.Context) {
				commandImport(c.GlobalBool("verbose"), c.BoolT("dc"))
			},
		},
		cli.Command{
			Name:  "export, e",
			Usage: "export from consul",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dc",
					Usage: "keep the dcs",
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
				commandExport(c.GlobalBool("verbose"), c.BoolT("dc"), c.BoolT("incversion"), c.String("prefix"))
			},
		},
	}

	app.Run(os.Args)
}

func initConsul(verbose bool) (*storage.Mssql, *consul.Consul) {
	config := configuration.GetConfig("config.json")
	config.Debug = verbose

	if config.Debug {
		log.Printf(config.PrintDebug())
	}

	// s := storage.Mssql{}
	conn := storage.OpenConnection(config)
	defer conn.Close()

	consul := consul.GetConsul(config)

	return conn, consul
}

func commandExport(verbose, keepDC, incversion bool, prefix string) {
	conn, consul := initConsul(verbose)

	dcs := consul.GetDCs()
	kvs := consul.GetKVs(prefix, dcs)
	conn.WriteKVs(kvs, keepDC)
}

func commandImport(verbose, keepDC bool) {
	conn, consul := initConsul(verbose)

	kvs, _ := conn.GetKVs()
	err := consul.WriteKVs(kvs, keepDC)
	if err != nil {
		log.Fatal(err)
	}
}

func commandValidate(file, validateHelpText string) {
	if file != "" {
		testConfiguration := configuration.GetConfig(file)
		os.Exit(testConfiguration.ValidateConfiguration())
	} else {
		fmt.Println(validateHelpText)
	}
}
