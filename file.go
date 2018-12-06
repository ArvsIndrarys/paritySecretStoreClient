package main

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
)

const (
	inputFolder     = "./input/"
	encryptedFolder = "./encrypted/"
	outputFolder    = "./output/"
	resultFolder    = "./result/"
)

// Path correspond to a file location, folder can be the folder's relative or full path
// filename is the name of the file in the folder
type Path struct {
	Folder   string
	FileName string
}

func (p Path) String() string {
	return buildString(p.Folder, "/", p.FileName)
}

func writeFile(content, out string) error {

	bytes := []byte(content)
	e := ioutil.WriteFile(out, bytes, 0644)
	return e
}

func resultToFile(result StoreResult, resultPath string) error {

	resultJSON, e := json.Marshal(result)
	if e != nil {
		return e
	}
	e = writeFile(string(resultJSON), resultPath)
	return e
}

func loadFile(path string) (string, error) {

	file, e := ioutil.ReadFile(path)
	if e != nil {
		return "", e
	}
	return string(file), nil
}

func loadFileAsHexString(path string) (string, error) {

	file, e := ioutil.ReadFile(path)
	if e != nil {
		return "", e
	}
	return hex.EncodeToString(file), nil
}

func loadInsertionResult(path string) (StoreResult, error) {

	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		return StoreResult{}, e
	}

	var res StoreResult
	e = json.Unmarshal(bytes, &res)
	if e != nil {
		return StoreResult{}, e
	}

	return res, nil
}
