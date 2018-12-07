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
	alicepwd             = "alicepwd"
	aliceAddr            = "0xb733e6b08e76559d5e9e55d5c56f10949c5017fd"
	signRawHashMethod    = "secretstore_signRawHash"
	generateDocKeyMethod = "secretstore_generateDocumentKey"
	encryptDocMethod     = "secretstore_encrypt"
	decryptDocMethod     = "secretstore_shadowDecrypt"
	shadow               = "shadow"
	version              = "2.0"
	id                   = 1
	treshold             = "2"
	plainDoc             = "fiftyShadesOfFailure"
)

var (
	BaseSecretStoreURL        net.URL
	BaseSecretStoreMethodsURL net.URL
	baseSignRawHashRequest    = net.Query{JSONRPCVersion: version, Method: signRawHashMethod, ID: id}
	baseGenerateDocKeyRequest = net.Query{JSONRPCVersion: version, Method: generateDocKeyMethod, ID: id}
	baseEncryptDocRequest     = net.Query{JSONRPCVersion: version, Method: encryptDocMethod, ID: id}
	baseDecryptDocRequest     = net.Query{JSONRPCVersion: version, Method: decryptDocMethod, ID: id}
)

func SignRandomHash() (string, error) {
	docID := randDocID()
	signedDocID, e := signRawHash(aliceAddr, alicepwd, docID)
	return signedDocID, e
}

func GenRandomKey() (string, error) {
	docID := randDocID()
	signed, e := signRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return "", e
	}

	serverKey, e := getServerKey(docID, signed)
	return serverKey, e
}

func ServerDocKeygen() (string, error) {
	docID := randDocID()
	signed, e := signRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return "", e
	}
	s, e := serverAndDocKeygen(docID, signed)
	return s, e
}

func InsertRandomDataInSecretStore() (string, error) {
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

	serverKey, e := getServerKey(docID, signedDocID)
	if e != nil {
		return fmt.Errorf("getDecryptionKey: %s", e)
	}

	encKey, e := generateDocumentKey(aliceAddr, alicepwd, serverKey)
	if e != nil {
		return fmt.Errorf("generateDocKey: %s", e)
	}

	hexData := hex.EncodeToString([]byte(plainDoc))

	fmt.Println("DataHex:", hexData)

	encDoc, e := encryptDocument(aliceAddr, alicepwd, encKey.EncryptedKey, hexData)
	if e != nil {
		return fmt.Errorf("encryptDocument: %s", e)
	}

	e = file.WriteEncryptedFile(encDoc, docID)
	if e != nil {
		return fmt.Errorf("writeFile: %s", e)
	}

	e = publishKey(docID, signedDocID, encKey.CommonPoint, encKey.EncryptedPoint)
	if e != nil {
		return fmt.Errorf("storeDoc: %s", e)
	}
	return nil
}

func signRawHash(address, password, documentID string) (string, error) {
	url := BaseSecretStoreMethodsURL

	request := baseSignRawHashRequest
	params := []string{address, password, strings.Add0x(documentID)}
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
	signedDocumentID := qr.Result

	return signedDocumentID, nil
}

func getServerKey(docID, signedDocID string) (string, error) {

	url := BaseSecretStoreURL
	path := strings.BuildString(shadow, "/", strings.Strip0x(docID), "/", strings.Strip0x(signedDocID), "/", treshold)
	url.Path = path

	resp, e := net.ExecutePost(url.String(), "")
	if e != nil {
		return "", e
	}

	return resp, nil
}

func generateDocumentKey(address, pwd, serverKey string) (parity.EncryptionKey, error) {

	url := BaseSecretStoreMethodsURL
	request := baseGenerateDocKeyRequest
	params := []string{address, pwd, serverKey}
	request.Params = params

	resp, e := net.ExecutePost(url.String(), request)
	if e != nil {
		return parity.EncryptionKey{}, e
	}

	var qr net.EncKeyQueryResult
	e = json.Unmarshal([]byte(resp), &qr)
	if e != nil {
		return parity.EncryptionKey{}, e
	}

	return qr.Result, nil
}

func encryptDocument(address, pwd, encKey, hexData string) (string, error) {

	url := BaseSecretStoreMethodsURL
	request := baseEncryptDocRequest
	params := []string{address, pwd, strings.Add0x(encKey), strings.Add0x(hexData)}
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
	encryptedDocument := qr.Result

	return encryptedDocument, nil
}

func publishKey(docID, signedDocID, commonPt, encryptedPt string) error {
	url := BaseSecretStoreURL
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

	url := BaseSecretStoreURL
	path := strings.BuildString(strings.Strip0x(docID), "/", strings.Strip0x(signedDocID)+"/"+treshold)
	url.Path = path
	s, e := net.ExecutePost(url.String(), "")
	return s, e
}
