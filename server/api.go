package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type MattermostAPI struct {
	ServerURL string
	AuthToken string
	Headers   map[string]string
}

func NewMattermostAPI(serverURL, authToken string) *MattermostAPI {
	return &MattermostAPI{
		ServerURL: serverURL,
		AuthToken: authToken,
		Headers: map[string]string{
			"Authorization": "Bearer " + authToken,
			"Content-Type":  "application/json",
		},
	}
}

func (api *MattermostAPI) endpointURL(path string) string {
	return fmt.Sprintf("%s/api/v4/%s", api.ServerURL, path)
}

func (api *MattermostAPI) check(response *http.Response) ([]byte, error) {
	if response != nil && (response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated) {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	} else if response != nil {
		fmt.Printf("Error %d: %s\n", response.StatusCode, response.Request.URL)
		return nil, fmt.Errorf("error with status code %d", response.StatusCode)
	} else {
		fmt.Println("Error - No response.")
		return nil, fmt.Errorf("no response from server")
	}
}

func (api *MattermostAPI) Get(path string) ([]byte, error) {
	url := api.endpointURL(path)
	req, _ := http.NewRequest("GET", url, nil)
	for key, value := range api.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return api.check(resp)
}

func (api *MattermostAPI) Post(path string, data interface{}) ([]byte, error) {
	url := api.endpointURL(path)
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	for key, value := range api.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return api.check(resp)
}

func (api *MattermostAPI) Delete(path string) ([]byte, error) {
	url := api.endpointURL(path)
	req, _ := http.NewRequest("DELETE", url, nil)
	for key, value := range api.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return api.check(resp)
}

func (api *MattermostAPI) Put(path string, data interface{}) ([]byte, error) {
	url := api.endpointURL(path)
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	for key, value := range api.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return api.check(resp)
}
