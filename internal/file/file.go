package file

import (
	"io/ioutil"
)

const (
	encryptedFolder = "./encrypted/"
)

func WriteEncryptedFile(content, out string) error {

	bytes := []byte(content)
	e := ioutil.WriteFile(encryptedFolder+out, bytes, 0644)
	return e
}

func LoadEncryptedFile(path string) (string, error) {

	file, e := ioutil.ReadFile(encryptedFolder + path)
	if e != nil {
		return "", e
	}
	return string(file), nil
}
