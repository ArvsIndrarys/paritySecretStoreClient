package core

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/config"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/strings"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/net"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/parity"
)

const (
	signRawHashMethod    = "secretstore_signRawHash"
	generateDocKeyMethod = "secretstore_generateDocumentKey"
	encryptDocMethod     = "secretstore_encrypt"
	decryptDocMethod     = "secretstore_shadowDecrypt"
)

var (
	baseSignRawHashRequest    = net.Query{JSONRPCVersion: version, Method: signRawHashMethod, ID: id}
	baseGenerateDocKeyRequest = net.Query{JSONRPCVersion: version, Method: generateDocKeyMethod, ID: id}
	baseEncryptDocRequest     = net.Query{JSONRPCVersion: version, Method: encryptDocMethod, ID: id}
	baseDecryptDocRequest     = net.Query{JSONRPCVersion: version, Method: decryptDocMethod, ID: id}
)

func SignRandomHash(address, password string) (string, error) {
	docID := randDocID()
	signedDocID, e := SignRawHash(address, password, docID)
	return signedDocID, e
}

func SignRawHash(address, password, documentID string) (string, error) {
	url := config.BaseParityRPCURL
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

func GenerateDocumentKey(address, pwd, serverKey string) (parity.EncryptionKey, error) {

	url := config.BaseParityRPCURL
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

func EncryptDocument(address, pwd, encKey, hexData string) (string, error) {

	url := config.BaseParityRPCURL
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

func DecryptDoc(address, pwd string, secret string, commonPoint string, shadows []string, encDoc string) (string, error) {

	decKey := parity.DecryptionKey{
		Secret:      secret,
		CommonPoint: commonPoint,
		Shadows:     shadows,
	}

	return decryptDoc(address, pwd, decKey, encDoc)
}

func decryptDoc(address, pwd string, decKey parity.DecryptionKey, encDoc string) (string, error) {
	url := config.BaseParityRPCURL
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
