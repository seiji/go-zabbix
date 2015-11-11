package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

const (
	MethodConfigurationExport = "configuration.export"
	MethodLogin = "user.login"
	MethodScreenGet = "screen.get"
	MethodTemplateGet = "template.get"
)

type ZabbixTemplate struct {
	Host string `json:"host"`
	Name string `json:"name"`
	TemplateID string `json:"templateid"`
	Description string `json:"description"`
}

type ZabbixScreen struct {
	Name string `json:"name"`
	ScrenID string `json:"screenid"`
	TemplateID string `json:"templateid"`
}

var (
	localhost, _ = os.Hostname()
	zabbixHost = os.Getenv("ZABBIX_HOST")
	zabbixUser = os.Getenv("ZABBIX_USER")
	zabbixPass = os.Getenv("ZABBIX_PASS")

	client *JsonRPCClient

	logLevel int

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

func setup(stage string) error {
	//logLevel = LOG_TRACE
	var err error
	client, err = NewClient(fmt.Sprintf("http://%s/zabbix/api_jsonrpc.php", zabbixHost), zabbixUser, zabbixPass)
	if err != nil {
		return err
	}

	return client.Login()
}

func templateExport(c *cli.Context) {
	err := setup(c.GlobalString("stage"))
	if err != nil {
		LogError(localhost, err)
		return
	}

	fmt.Printf("client = %+v\n", client)

	templates, err := client.TemplateList()
	if err != nil {
		LogError(localhost, err)
	}

	client.TemplateExport(c.String("dir"), "xml", templates...)
}

func screenExport(c *cli.Context) {
	err := setup(c.GlobalString("stage"))
	if err != nil {
		LogError(localhost, err)
		return
	}

	fmt.Printf("client = %+v\n", client)

	screens, err := client.ScreenList()
	if err != nil {
		LogError(localhost, err)
	}
	fmt.Printf("screens = %+v\n", screens)

	client.ScreenExport(c.String("dir"), "xml", screens...)
}

