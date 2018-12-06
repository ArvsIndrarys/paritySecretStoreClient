package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func decryptViaStore(docID string) (string, error) {

	data, e := loadInsertionResult(resultFolder + docID)
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

	encDoc, e := loadFile(data.EncryptedDocumentPath)

	plainDoc, e := decryptDoc(bobAddr, bobpwd, decKeys, encDoc)
	if e != nil {
		return "", fmt.Errorf("decryptDoc: %s", e)
	}

	e = writeFile(plainDoc, outputFolder+docID)
	if e != nil {
		return "", fmt.Errorf("writeFile: %s", e)
	}

	return plainDoc, nil
}

func getDecryptionKeys(docID, signedDocID string) (DecryptionKey, error) {
	url := baseSecretStoreURL
	path := buildString(shadow, "/", strip0x(docID), "/", strip0x(signedDocID))
	url.Path = path

	resp, e := ExecuteGet(url.String())
	if e != nil {
		return DecryptionKey{}, e
	}

	var decKeys DecryptionKey
	e = json.Unmarshal([]byte(resp), &decKeys)
	if e != nil {
		return DecryptionKey{}, e
	}
	return decKeys, nil
}

func decryptDoc(address, pwd string, decKey DecryptionKey, encDoc string) (string, error) {
	url := baseSecretStoreMethodsURL
	request := baseDecryptDocRequest
	params := []string{address, pwd,
		decKey.Secret, decKey.CommonPoint, decKey.GetShadowsString(),
		encDoc}

	request.Params = params

	resp, e := ExecutePost(url.String(), request)
	if e != nil {
		return "", e
	}

	var qr QueryResult
	e = json.Unmarshal([]byte(resp), &qr)
	if e != nil {
		return "", e
	}
	decryptedDocument := strip0x(qr.Result)

	plainBytes, e := hex.DecodeString(decryptedDocument)
	if e != nil {
		return "", e
	}

	return string(plainBytes), nil
}
