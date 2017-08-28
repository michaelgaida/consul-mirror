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
		cli.StringFlag{
			Name:  "validate",
			Value: "",
			Usage: "configuration file to be validated",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Switch on the verbose mode",
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "import",
			Usage: "import from consul",
		},
		cli.Command{
			Name:  "export",
			Usage: "export from consul",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dc",
					Usage: "keep the dcs",
				},
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.String("validate") != "" {
			testConfiguration := configuration.GetConfig(c.String("validate"))
			os.Exit(testConfiguration.ValidateConfiguration())
		}
		if (!(c.Bool("import"))) && (!(c.Bool("export"))) {
			fmt.Println("Nothing to do")
			os.Exit(0)
		}
		config := configuration.GetConfig("config.json")
		if c.Bool("verbose") {
			config.Debug = true
		}
		if config.Debug {
			log.Printf(config.PrintDebug())
		}

		// s := storage.Mssql{}
		conn := storage.OpenConnection(config)
		// defer conn.Close()

		consul := consul.GetConsul(config)

		dcs := consul.GetDCs()

		if c.Bool("export") {
			kvs := consul.GetKVs("", dcs)
			if kvs == nil {
			}
			conn.WriteKVs(kvs)
		}

		if c.Bool("import") {
			kvs, _ := conn.GetKVs()
			err := consul.WriteKVs(kvs, c.Bool("dc"))
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	}

	app.Run(os.Args)
}
