package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// QueryResult is a base Parity response
type QueryResult struct {
	JSONRPCVersion string `json:"jsonrpc"`
	Result         string `json:"result"`
	ID             int    `json:"id"`
}

// URL is a decomposition of an URL
type URL struct {
	BaseURL string
	Port    string
	Path    string
}

func (u URL) String() string {
	return buildString(u.BaseURL, ":", u.Port, "/", u.Path)
}

// ExecutePost sends a HTTP POST request with a JSON object, if obj is different from an empty string
// returns the body of the result as a string
func ExecutePost(urlPath string, obj interface{}) (string, error) {

	var e error
	var jsonObj []byte

	if obj != "" {
		jsonObj, e = json.Marshal(obj)
		if e != nil {
			log.Println(e)
		}
	}
	jsonObj = formatJSON(jsonObj)

	buf := bytes.NewBuffer(jsonObj)

	req, e := http.NewRequest(http.MethodPost, urlPath, buf)
	if e != nil {
		return "", e
	}
	req.Header.Set("Content-Type", "application/json")

	resp, e := sendRequest(req)
	if e != nil {
		return "", e
	}
	return resp, nil
}

// ExecuteGet sens a HTTP GET request
// returns the body of the result as a string
func ExecuteGet(url string) (string, error) {

	req, e := http.NewRequest(http.MethodGet, url, nil)
	resp, e := sendRequest(req)
	if e != nil {
		return "", e
	}
	return resp, e
}

// sendRequest adds an additional error check on the response returned by Parity Secret Sharing
func sendRequest(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body := bodyToString(resp.Body)
	if resp.StatusCode != 200 || strings.Contains(body, "error") {
		return "", fmt.Errorf("Error calling %s: \n[code: %s message: %v]", req.URL, resp.Status, body)
	}

	return body, nil
}

func bodyToString(b io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)
	return strings.Trim(buf.String(), "\"")
}

// resolves Shadow ([]String) problem
func formatJSON(in []byte) []byte {
	outStr := strings.Replace(string(in), "\\", "", -1)
	outStr = strings.Replace(outStr, "\"[", "[", -1)
	outStr = strings.Replace(outStr, "]\"", "]", -1)
	return []byte(outStr)
}
