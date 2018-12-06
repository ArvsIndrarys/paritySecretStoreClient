package main

import (
	"encoding/json"
	"fmt"
)

var (
	baseSecretStoreURL        URL
	baseSecretStoreMethodsURL URL
	baseSignRawHashRequest    = Query{JSONRPCVersion: version, Method: signRawHashMethod, ID: id}
	baseGenerateDocKeyRequest = Query{JSONRPCVersion: version, Method: generateDocKeyMethod, ID: id}
	baseEncryptDocRequest     = Query{JSONRPCVersion: version, Method: encryptDocMethod, ID: id}
	baseDecryptDocRequest     = Query{JSONRPCVersion: version, Method: decryptDocMethod, ID: id}
)

// StoreResult are the informations to be shared to someone to be able to decrypt the file
type StoreResult struct {
	DocumentID            string `json:"document_id"`
	SignedDocumentID      string `json:"signed_document_id"`
	EncryptedDocumentPath string `json:"encrypted_document_path"`
}

func (r StoreResult) String() string {
	return buildString("Document Key Id: ", r.DocumentID, "\nSigned Document Key Id: ",
		r.SignedDocumentID, "\nResulting encrypted document location: ", r.EncryptedDocumentPath)
}

func signRandomHash() (string, error) {
	docID := randDocID()
	signedDocID, e := signRawHash(aliceAddr, alicepwd, docID)
	return signedDocID, e
}

func genRandomKey() (string, error) {
	docID := randDocID()
	signed, e := signRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return "", e
	}

	serverKey, e := getServerKey(docID, signed)
	return serverKey, e
}

func serverDocKeygen() (string, error) {
	docID := randDocID()
	signed, e := signRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return "", e
	}
	s, e := serverAndDocKeygen(docID, signed)
	return s, e
}

func insertRandomDataInSecretStore() (string, error) {
	docID := randDocID()
	e := insertDataInSecretStore(docID)
	if e != nil {
		return "", e
	}
	return docID, nil
}

func insertDataInSecretStore(docID string) error {

	signedDocID, e := signRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return fmt.Errorf("signRawHash: %s", e)
	}

	fmt.Println("DocID", docID, "\nsignedDocID", signedDocID)

	serverKey, e := getServerKey(docID, signedDocID)
	if e != nil {
		return fmt.Errorf("getDecryptionKey: %s", e)
	}

	encKey, e := generateDocumentKey(aliceAddr, alicepwd, serverKey)
	if e != nil {
		return fmt.Errorf("generateDocKey: %s", e)
	}

	hexData, e := loadFileAsHexString(inputFolder + docID)
	if e != nil {
		return fmt.Errorf("loadFileAsHexString: %s", e)
	}

	encDoc, e := encryptDocument(aliceAddr, alicepwd, encKey.EncryptedKey, hexData)
	if e != nil {
		return fmt.Errorf("encryptDocument: %s", e)
	}

	e = writeFile(encDoc, encryptedFolder+docID)
	if e != nil {
		return fmt.Errorf("writeFile: %s", e)
	}

	e = publishKey(docID, signedDocID, encKey.CommonPoint, encKey.EncryptedPoint)
	if e != nil {
		return fmt.Errorf("storeDoc: %s", e)
	}

	result := StoreResult{DocumentID: docID, SignedDocumentID: signedDocID, EncryptedDocumentPath: encryptedFolder + docID}
	e = resultToFile(result, resultFolder+docID)
	if e != nil {
		return fmt.Errorf("resultToFile: %s", e)
	}
	return nil
}

func signRawHash(address, password, documentID string) (string, error) {
	url := baseSecretStoreMethodsURL

	request := baseSignRawHashRequest
	params := []string{address, password, add0x(documentID)}
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
	signedDocumentID := qr.Result

	return signedDocumentID, nil
}

func getServerKey(docID, signedDocID string) (string, error) {

	url := baseSecretStoreURL
	path := buildString(shadow, "/", strip0x(docID), "/", strip0x(signedDocID), "/", treshold)
	url.Path = path

	resp, e := ExecutePost(url.String(), "")
	if e != nil {
		return "", e
	}

	return resp, nil
}

func generateDocumentKey(address, pwd, serverKey string) (EncryptionKey, error) {

	url := baseSecretStoreMethodsURL
	request := baseGenerateDocKeyRequest
	params := []string{address, pwd, serverKey}
	request.Params = params

	resp, e := ExecutePost(url.String(), request)
	if e != nil {
		return EncryptionKey{}, e
	}

	var qr EncKeyQueryResult
	e = json.Unmarshal([]byte(resp), &qr)
	if e != nil {
		return EncryptionKey{}, e
	}

	return qr.Result, nil
}

func encryptDocument(address, pwd, encKey, hexData string) (string, error) {

	url := baseSecretStoreMethodsURL
	request := baseEncryptDocRequest
	params := []string{address, pwd, add0x(encKey), add0x(hexData)}
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
	encryptedDocument := qr.Result

	return encryptedDocument, nil
}

func publishKey(docID, signedDocID, commonPt, encryptedPt string) error {
	url := baseSecretStoreURL
	path := buildString(shadow, "/", strip0x(docID), "/", strip0x(signedDocID), "/",
		strip0x(commonPt), "/", strip0x(encryptedPt))
	url.Path = path

	_, e := ExecutePost(url.String(), "")
	if e != nil {
		return e
	}
	return nil
}

func serverAndDocKeygen(docID, signedDocID string) (string, error) {

	url := baseSecretStoreURL
	path := buildString(strip0x(docID), "/", strip0x(signedDocID)+"/"+treshold)
	url.Path = path
	s, e := ExecutePost(url.String(), "")
	return s, e
}
