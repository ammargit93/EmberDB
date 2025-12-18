package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Data struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Client struct {
	Addr string
}

func Connect(addr string) Client {
	var client Client
	client.Addr = addr
	return client
}

func (client Client) SetKey(key string, value any) (Data, error) {
	var data Data
	data.Key = key
	data.Value = value
	ct := &http.Client{}

	dataBytes, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", client.Addr+"/set", bytes.NewBuffer(dataBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := ct.Do(req)
	if err != nil {
		return Data{}, err
	}
	defer resp.Body.Close()
	respContent, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respContent))
	var response Data
	json.Unmarshal(respContent, &response)
	return response, err
}
