package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewRestClient(serverURL, authToken string, headers map[string]string) *RestClient {
	return &RestClient{
		ServerURL: serverURL,
		AuthToken: authToken,
		Headers:   headers,
		Client:    &http.Client{},
	}
}

type RestClient struct {
	ServerURL string
	AuthToken string
	Headers   map[string]string
	Client    *http.Client
}

func (r *RestClient) endpointURL(path string) string {
	return fmt.Sprintf("%s/api/v4/%s", r.ServerURL, path)
}

func (r *RestClient) check(response *http.Response) ([]byte, error) {
	if response != nil && (response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated) {
		body, err := io.ReadAll(response.Body)
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

func (r *RestClient) Get(path string) ([]byte, error) {
	url := r.endpointURL(path)
	req, _ := http.NewRequest("GET", url, nil)
	for key, value := range r.Headers {
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

	return r.check(resp)
}

func (r *RestClient) Post(path string, data interface{}) ([]byte, error) {
	url := r.endpointURL(path)
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	for key, value := range r.Headers {
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

	return r.check(resp)
}

func (r *RestClient) Delete(path string) ([]byte, error) {
	url := r.endpointURL(path)
	req, _ := http.NewRequest("DELETE", url, nil)
	for key, value := range r.Headers {
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

	return r.check(resp)
}

func (r *RestClient) Put(path string, data interface{}) ([]byte, error) {
	url := r.endpointURL(path)
	jsonData, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	for key, value := range r.Headers {
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

	return r.check(resp)
}
