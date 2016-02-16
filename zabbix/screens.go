package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	MethodConfigurationExport = "configuration.export"
	MethodScreenGet           = "screen.get"
)

type ZabbixScreen struct {
	Name       string `json:"name"`
	ScrenID    string `json:"screenid"`
	TemplateID string `json:"templateid"`
}

type jsonRPCResponse struct {
	Jsonrpc string       `json:"jsonrpc"`
	Error   JsonRPCError `json:"error"`
	Result  interface{}  `json:"result"`
	Id      int          `json:"id"`
}

func (c *JsonRPCClient) ScreenList() ([]ZabbixScreen, error) {
	var res JsonRPCResponseScreens

	data := map[string]interface{}{
		"output": "extend",
	}

	err := c.request(MethodScreenGet, data, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *JsonRPCClient) ScreenExport(dirPath, format string, templates ...ZabbixScreen) error {

	for _, v := range templates {
		var res JsonRPCResponseGeneral

		data := map[string]interface{}{
			"options": map[string][]string{"screens": []string{v.ScrenID}},
			"format":  format,
		}

		err := c.request(MethodConfigurationExport, data, &res)
		if err != nil {
			return err
		}

		fullPath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", v.Name, format))
		fout, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer fout.Close()
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, []byte(res.Result), "", "    ")
		if err != nil {
			return err
		}

		fout.WriteString(string(prettyJSON.Bytes()))
		fmt.Printf("Exported %+v\n", fullPath)
	}

	return nil
}
