package file

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/strings"
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
	return strings.BuildString(p.Folder, "/", p.FileName)
}

// StoreResult are the informations to be shared to someone to be able to decrypt the file
type StoreResult struct {
	DocumentID            string `json:"document_id"`
	SignedDocumentID      string `json:"signed_document_id"`
	EncryptedDocumentPath string `json:"encrypted_document_path"`
}

func (r StoreResult) String() string {
	return strings.BuildString("Document Key Id: ", r.DocumentID,
		"\nSigned Document Key Id: ", r.SignedDocumentID,
		"\nResulting encrypted document location: ", r.EncryptedDocumentPath)
}

func WriteFile(content, out string) error {

	bytes := []byte(content)
	e := ioutil.WriteFile(outputFolder+out, bytes, 0644)
	return e
}

func ResultToFile(result StoreResult, resultPath string) error {

	resultJSON, e := json.Marshal(result)
	if e != nil {
		return e
	}
	e = WriteFile(string(resultJSON), resultFolder+resultPath)
	return e
}

func LoadFile(path string) (string, error) {

	file, e := ioutil.ReadFile(path)
	if e != nil {
		return "", e
	}
	return string(file), nil
}

func LoadFileAsHexString(path string) (string, error) {

	file, e := ioutil.ReadFile(inputFolder + path)
	if e != nil {
		return "", e
	}
	return hex.EncodeToString(file), nil
}

func LoadInsertionResult(path string) (StoreResult, error) {

	bytes, e := ioutil.ReadFile(resultFolder + path)
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
