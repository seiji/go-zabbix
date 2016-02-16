package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/seiji/go-log/log"
	"github.com/seiji/go-zabbix/zabbix"
)

var (
	version      = "HEAD"
	curDir, _    = os.Getwd()
	localhost, _ = os.Hostname()
	zabbixHost   = os.Getenv("ZABBIX_HOST")
	zabbixUser   = os.Getenv("ZABBIX_USER")
	zabbixPass   = os.Getenv("ZABBIX_PASS")

	env map[string]string
)

func init() {
	if zabbixHost == "" {
		zabbixHost = "127.0.0.1"
	}
	if zabbixUser == "" {
		zabbixUser = "admin"
	}
	if zabbixPass == "" {
		zabbixPass = "zabbix"
	}
}

func main() {
	newApp().Run(os.Args)
}

func newApp() (app *cli.App) {
	app = cli.NewApp()
	app.Name = "zabbix-api"
	app.Usage = "zabbix api client"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:    "template_export",
			Aliases: []string{"t"},
			Usage:   "export templates",
			Action:  templateExport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Value: curDir,
					Usage: "dir to export",
				},
			},
		},
		{
			Name:    "screen_export",
			Aliases: []string{"s"},
			Usage:   "export screens",
			Action:  screenExport,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir, d",
					Value: curDir,
					Usage: "dir to export",
				},
			},
		},
	}
	return
}

func templateExport(c *cli.Context) {
	client := zabbix.NewClient(zabbixHost, zabbixUser, zabbixPass)
	client.Login()

	templates, err := client.TemplateList()
	if err != nil {
		log.LogError(err)
	}

	fmt.Printf("templates = %+v\n", templates)
	client.TemplateExport(c.String("dir"), "json", templates...)
}

func screenExport(c *cli.Context) {
	client := zabbix.NewClient(zabbixHost, zabbixUser, zabbixPass)
	client.Login()

	screens, err := client.ScreenList()
	if err != nil {
		log.LogError(err)
	}
	fmt.Printf("screens = %+v\n", screens)

	client.ScreenExport(c.String("dir"), "json", screens...)
}
