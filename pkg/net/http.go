package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/parity"
)

// Query is a base Parity query
type Query struct {
	JSONRPCVersion string   `json:"jsonrpc"`
	Method         string   `json:"method"`
	Params         []string `json:"params"`
	ID             int      `json:"id"`
}

// EncKeyQueryResult is a Parity response containing an encryption Key
type EncKeyQueryResult struct {
	JSONRPCVersion string               `json:"jsonrpc"`
	Result         parity.EncryptionKey `json:"result"`
	ID             int                  `json:"id"`
}

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
	var b strings.Builder
	fmt.Fprintf(&b, "%s:%s/%s", u.BaseURL, u.Port, u.Path)
	return b.String()
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

// resolves Shadow ([]String) problem the dirty way
func formatJSON(in []byte) []byte {
	outStr := strings.Replace(string(in), "\\", "", -1)
	outStr = strings.Replace(outStr, "\"[", "[", -1)
	outStr = strings.Replace(outStr, "]\"", "]", -1)
	return []byte(outStr)
}
