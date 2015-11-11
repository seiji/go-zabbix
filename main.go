package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var (
	version string = "HEAD"
	curDir, _ = os.Getwd()
)

func main() {
	newApp().Run(os.Args)
}

func newApp() (app *cli.App) {
	app = cli.NewApp()
	app.Name = "zabbix-api"
	app.Usage = "zabbix api client"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "stage, s",
			Value: "prod",
			Usage: "stage for deploy",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "template_export",
			Aliases: []string{"t"},
			Usage:	 "export templates",
			Action:  templateExport,
			Flags:   []cli.Flag{
				cli.StringFlag{
					Name: "dir, d",
					Value: curDir,
					Usage: "dir to export",
				},
			},
		},
		{
			Name:    "screen_export",
			Aliases: []string{"s"},
			Usage:	 "export screens",
			Action:  screenExport,
			Flags:   []cli.Flag{
				cli.StringFlag{
					Name: "dir, d",
					Value: curDir,
					Usage: "dir to export",
				},
			},
		},
	}
	return
}
