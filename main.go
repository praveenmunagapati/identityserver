package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Identity server"
	app.Version = "0.1-Dev"

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	var debugLogging bool
	var bindAddress string
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			Name:        "bind, b",
			Usage:       "Bind address",
			Value:       ":8080",
			Destination: &bindAddress,
		},
	}
	app.Before = func(c *cli.Context) error {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
			log.Debug(app.Name, "-", app.Version)
		}
		return nil
	}

	app.Action = func(c *cli.Context) {

		//r, err := server.GetRouter()

		//if err != nil {
		//	log.Error(err)
		//	os.Exit(1)
		//}

		//r.Run(bindAddress)
	}

	app.Run(os.Args)
}
