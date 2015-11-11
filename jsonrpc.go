package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type JsonRPCClient struct {
	address string
	user string
	pass string
	id  int
	auth string
}

type JsonRPCRequest struct {
	Jsonrpc string			`json:"jsonrpc"`
	Method  string			`json:"method"`
	Params  interface{}		`json:"params"`

	Auth	string			`json:"auth,omitempty"`
	Id		int				`json:"id"`
}

type omit *struct {}

type JsonRPCError struct {
	Code    int				`json:"code"`
	Message string			`json:"message"`
	Data    string			`json:"data"`
}

func (e *JsonRPCError) Error() string {
	return e.Data
}

type JsonRPCResponse interface { GetError() *JsonRPCError}

type jsonRPCResponse struct {
	Jsonrpc string			`json:"jsonrpc"`
	Error   JsonRPCError	`json:"error"`
	Result  interface{}		`json:"result"`
	Id      int				`json:"id"`
}

func (r *jsonRPCResponse) GetError() *JsonRPCError {
	return &r.Error
}

type JsonRPCResponseGeneral struct {
	*jsonRPCResponse
	Result  string			`json:"result"`
}

type JsonRPCResponseTemplates struct {
	*jsonRPCResponse
	Result []ZabbixTemplate `json:"result"`
}

type JsonRPCResponseScreens struct {
	*jsonRPCResponse
	Result []ZabbixScreen `json:"result"`
}

func (c *JsonRPCClient) request(method string, data map[string]interface{}, result JsonRPCResponse) error {
	c.id = c.id + 1
	jsonobj := JsonRPCRequest{"2.0", method, data, c.auth, c.id}
	encoded, err := json.Marshal(jsonobj)

	request, err := http.NewRequest("POST", c.address, bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json-rpc")

	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var buf bytes.Buffer

	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(buf.Bytes(), result)
	if result.GetError().Code != 0 {
		return  result.GetError()
	}

	return nil
}

func (c *JsonRPCClient) Login() error{
	var res JsonRPCResponseGeneral

	data :=  map[string]interface{} {
		"user": c.user,
		"password": c.pass,
	}

	err :=  c.request(MethodLogin, data, &res)
	if err != nil {
		return err
	}
	c.auth = res.Result

	return nil
}

func (c *JsonRPCClient) TemplateList() ([]ZabbixTemplate, error){
	var res JsonRPCResponseTemplates

	data :=  map[string]interface{} {
		"output": "extend",
	}

	err :=  c.request(MethodTemplateGet, data, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *JsonRPCClient) TemplateExport(dirPath, format string, templates ...ZabbixTemplate) error {
	for _, v := range templates {
		var res JsonRPCResponseGeneral

		data := map[string]interface{} {
			"options": map[string][]string{"templates": []string{v.TemplateID}},
			"format" : format,
		}

		err :=  c.request(MethodConfigurationExport, data, &res)
		if err != nil {
			return err
		}

		fullPath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", v.Name, format))
		fout, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer fout.Close()
		fout.WriteString(res.Result)
		fmt.Printf("Exported %+v\n", fullPath)
	}
	return nil
}

func (c *JsonRPCClient) ScreenList() ([]ZabbixScreen, error){
	var res JsonRPCResponseScreens

	data :=  map[string]interface{} {
		"output": "extend",
	}

	err :=  c.request(MethodScreenGet, data, &res)
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *JsonRPCClient) ScreenExport(dirPath, format string, templates ...ZabbixScreen) error {

	for _, v := range templates {
		var res JsonRPCResponseGeneral

		data := map[string]interface{} {
			"options": map[string][]string{"screens": []string{v.ScrenID}},
			"format" : format,
		}

		err :=  c.request(MethodConfigurationExport, data, &res)
		if err != nil {
			return err
		}

		fullPath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", v.Name, format))
		fout, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		defer fout.Close()
		fout.WriteString(res.Result)
		fmt.Printf("Exported %+v\n", fullPath)
	}

	return nil
}

func NewClient (address, user, pass string) (*JsonRPCClient, error) {
	return &JsonRPCClient{address, user, pass, 0, "",}, nil
}

