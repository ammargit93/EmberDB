package sdk

import (
	"bytes"
	"encoding/json"
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

type MGetResponse struct {
	Values map[string]int `json:"values"`
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
	var response Data
	json.Unmarshal(respContent, &response)
	return response, err
}

func (client Client) GetKey(key string) (Data, error) {

	ct := &http.Client{}

	req, _ := http.NewRequest("GET", client.Addr+"/get/"+key, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := ct.Do(req)
	if err != nil {
		return Data{}, err
	}
	defer resp.Body.Close()
	respContent, _ := io.ReadAll(resp.Body)
	var response Data
	json.Unmarshal(respContent, &response)
	return response, err
}

func (client Client) MSetKey(pairs map[string]any) (map[string]any, error) {
	ct := &http.Client{}
	dataBytes, _ := json.Marshal(pairs)
	req, _ := http.NewRequest("POST", client.Addr+"/mset", bytes.NewBuffer(dataBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := ct.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respContent, _ := io.ReadAll(resp.Body)
	var response map[string]any
	json.Unmarshal(respContent, &response)
	return response, err
}
func (c Client) MGetKey(keys []string) (map[string]int, error) {
	ct := &http.Client{}

	dataBytes, _ := json.Marshal(keys)
	req, _ := http.NewRequest("POST", c.Addr+"/mget", bytes.NewBuffer(dataBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := ct.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respContent, _ := io.ReadAll(resp.Body)

	var response struct {
		Values map[string]int `json:"values"`
	}

	if err := json.Unmarshal(respContent, &response); err != nil {
		return nil, err
	}

	return response.Values, nil
}
