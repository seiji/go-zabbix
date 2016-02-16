package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	MethodLogin = "user.login"
)

type JsonRPCClient struct {
	address string
	user    string
	pass    string
	id      int
	auth    string
}

type JsonRPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`

	Auth string `json:"auth,omitempty"`
	Id   int    `json:"id"`
}

type omit *struct{}

type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *JsonRPCError) Error() string {
	return e.Data
}

type JsonRPCResponse interface {
	GetError() *JsonRPCError
}

func (r *jsonRPCResponse) GetError() *JsonRPCError {
	return &r.Error
}

type JsonRPCResponseGeneral struct {
	*jsonRPCResponse
	Result string `json:"result"`
}

type JsonRPCResponseScreens struct {
	*jsonRPCResponse
	Result []ZabbixScreen `json:"result"`
}

func NewClient(host, user, pass string) *JsonRPCClient {
	return &JsonRPCClient{fmt.Sprintf("http://%s/zabbix/api_jsonrpc.php", host), user, pass, 0, ""}
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
		return result.GetError()
	}

	return nil
}

func (c *JsonRPCClient) Login() error {
	var res JsonRPCResponseGeneral

	data := map[string]interface{}{
		"user":     c.user,
		"password": c.pass,
	}

	err := c.request(MethodLogin, data, &res)
	if err != nil {
		return err
	}
	c.auth = res.Result

	return nil
}
