package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	MethodTemplateGet = "template.get"
)

type ZabbixTemplate struct {
	Host        string `json:"host"`
	Name        string `json:"name"`
	TemplateID  string `json:"templateid"`
	Description string `json:"description"`
}

type JsonRPCResponseTemplates struct {
	*jsonRPCResponse
	Result []ZabbixTemplate `json:"result"`
}

func (c *JsonRPCClient) TemplateList() ([]ZabbixTemplate, error) {
	var res JsonRPCResponseTemplates

	data := map[string]interface{}{
		"output": "extend",
	}

	err := c.request(MethodTemplateGet, data, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *JsonRPCClient) TemplateExport(dirPath, format string, templates ...ZabbixTemplate) error {
	for _, v := range templates {
		var res JsonRPCResponseGeneral

		data := map[string]interface{}{
			"options": map[string][]string{"templates": []string{v.TemplateID}},
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
