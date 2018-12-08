package core

import (
	"encoding/json"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/config"
	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/strings"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/net"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/parity"
)

const (
	shadow   = "shadow"
	version  = "2.0"
	id       = 1
	treshold = "2"
	plainDoc = "fiftyShadesOfFailure"
)

func GetServerKey(docID, signedDocID string) (string, error) {

	url := config.BaseSecretStoreURL
	path := strings.BuildString(shadow, "/", strings.Strip0x(docID), "/", strings.Strip0x(signedDocID), "/", treshold)
	url.Path = path

	resp, e := net.ExecutePost(url.String(), "")
	if e != nil {
		return "", e
	}

	return resp, nil
}

func StoreKey(docID, signedDocID, commonPt, encryptedPt string) error {
	url := config.BaseSecretStoreURL
	path := strings.BuildString(shadow, "/", strings.Strip0x(docID), "/", strings.Strip0x(signedDocID), "/",
		strings.Strip0x(commonPt), "/", strings.Strip0x(encryptedPt))
	url.Path = path

	_, e := net.ExecutePost(url.String(), "")
	if e != nil {
		return e
	}
	return nil
}

func serverAndDocKeygen(docID, signedDocID string) (string, error) {

	url := config.BaseSecretStoreURL
	path := strings.BuildString(strings.Strip0x(docID), "/", strings.Strip0x(signedDocID)+"/"+treshold)
	url.Path = path
	s, e := net.ExecutePost(url.String(), "")
	return s, e
}

func getDecryptionKeys(docID, signedDocID string) (parity.DecryptionKey, error) {
	url := config.BaseSecretStoreURL
	path := strings.BuildString(shadow, "/", strings.Strip0x(docID), "/", strings.Strip0x(signedDocID))
	url.Path = path

	resp, e := net.ExecuteGet(url.String())
	if e != nil {
		return parity.DecryptionKey{}, e
	}

	var decKeys parity.DecryptionKey
	e = json.Unmarshal([]byte(resp), &decKeys)
	if e != nil {
		return parity.DecryptionKey{}, e
	}
	return decKeys, nil
}
