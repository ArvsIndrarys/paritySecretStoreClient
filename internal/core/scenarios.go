package core

import (
	"encoding/hex"
	"fmt"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/file"
)

const (
	alicepwd  = "alicepwd"
	aliceAddr = "0xb733e6b08e76559d5e9e55d5c56f10949c5017fd"
	bobpwd    = "bobpwd"
	bobAddr   = "0x7b51be7a60c022fba408e6a82a9bf69e71a18528"
)

func InsertRandomDataInSecretStore() (string, error) {
	docID := randDocID()
	e := InsertDataInSecretStore(docID)
	if e != nil {
		return "", e
	}
	return docID, nil
}

func InsertDataInSecretStore(docID string) error {

	signedDocID, e := SignRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return fmt.Errorf("signRawHash: %s", e)
	}

	serverKey, e := GetServerKey(docID, signedDocID)
	if e != nil {
		return fmt.Errorf("getDecryptionKey: %s", e)
	}

	encKey, e := GenerateDocumentKey(aliceAddr, alicepwd, serverKey)
	if e != nil {
		return fmt.Errorf("generateDocKey: %s", e)
	}

	hexData := hex.EncodeToString([]byte(plainDoc))

	encDoc, e := EncryptDocument(aliceAddr, alicepwd, encKey.EncryptedKey, hexData)
	if e != nil {
		return fmt.Errorf("encryptDocument: %s", e)
	}

	e = file.WriteEncryptedFile(encDoc, docID)
	if e != nil {
		return fmt.Errorf("writeFile: %s", e)
	}

	e = StoreKey(docID, signedDocID, encKey.CommonPoint, encKey.EncryptedPoint)
	if e != nil {
		return fmt.Errorf("storeDoc: %s", e)
	}
	return nil
}

func GenRandomKey() (string, error) {
	docID := randDocID()
	signed, e := SignRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return "", e
	}

	serverKey, e := GetServerKey(docID, signed)
	return serverKey, e
}

func RandDocAndKeygen() (string, error) {
	docID := randDocID()
	signed, e := SignRawHash(aliceAddr, alicepwd, docID)
	if e != nil {
		return "", e
	}
	s, e := serverAndDocKeygen(docID, signed)
	return s, e
}

func DecryptViaStore(docID string) (string, error) {

	signedDocID, e := SignRawHash(bobAddr, bobpwd, docID)
	if e != nil {
		return "", fmt.Errorf("signRawHash: %s", e)
	}

	decKeys, e := getDecryptionKeys(docID, signedDocID)
	if e != nil {
		return "", fmt.Errorf("getDecryptionKeys: %s", e)
	}

	encDoc, e := file.LoadEncryptedFile(docID)

	plainDoc, e := decryptDoc(bobAddr, bobpwd, decKeys, encDoc)
	if e != nil {
		return "", fmt.Errorf("decryptDoc: %s", e)
	}

	return plainDoc, nil
}
