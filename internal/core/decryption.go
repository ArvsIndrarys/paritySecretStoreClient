package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/file"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/strings"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/net"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/parity"
)

const (
	bobpwd  = "bobpwd"
	bobAddr = "0x7b51be7a60c022fba408e6a82a9bf69e71a18528"
)

func DecryptViaStore(docID string) (string, error) {

	data, e := file.LoadInsertionResult(docID)
	if e != nil {
		return "", fmt.Errorf("loadResult: %s", e)
	}

	signedDocID, e := signRawHash(bobAddr, bobpwd, data.DocumentID)
	if e != nil {
		return "", fmt.Errorf("signRawHash: %s", e)
	}

	decKeys, e := getDecryptionKeys(data.DocumentID, signedDocID)
	if e != nil {
		return "", fmt.Errorf("getDecryptionKeys: %s", e)
	}

	encDoc, e := file.LoadFile(data.EncryptedDocumentPath)

	plainDoc, e := decryptDoc(bobAddr, bobpwd, decKeys, encDoc)
	if e != nil {
		return "", fmt.Errorf("decryptDoc: %s", e)
	}

	e = file.WriteFile(plainDoc, docID)
	if e != nil {
		return "", fmt.Errorf("writeFile: %s", e)
	}

	return plainDoc, nil
}

func getDecryptionKeys(docID, signedDocID string) (parity.DecryptionKey, error) {
	url := BaseSecretStoreURL
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

func decryptDoc(address, pwd string, decKey parity.DecryptionKey, encDoc string) (string, error) {
	url := BaseSecretStoreMethodsURL
	request := baseDecryptDocRequest
	params := []string{address, pwd,
		decKey.Secret, decKey.CommonPoint, decKey.GetShadowsString(),
		encDoc}

	request.Params = params

	resp, e := net.ExecutePost(url.String(), request)
	if e != nil {
		return "", e
	}

	var qr net.QueryResult
	e = json.Unmarshal([]byte(resp), &qr)
	if e != nil {
		return "", e
	}
	decryptedDocument := strings.Strip0x(qr.Result)

	plainBytes, e := hex.DecodeString(decryptedDocument)
	if e != nil {
		return "", e
	}

	return string(plainBytes), nil
}
