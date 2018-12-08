package server

import (
	"log"
	"net/http"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/core"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/parity"

	"github.com/gin-gonic/gin"
)

// ----- STRUCTURES -----
// -- Requests --

type signRawHashRequest struct {
	parity.Credentials
	DocumentID string `json:"docId"`
}

func (s signRawHashRequest) isCorrect() bool {
	return s.Credentials.IsCorrect() && s.DocumentID != ""
}

type docKeygenRequest struct {
	parity.Credentials
	ServerPubKey string `json:"pubKey"`
}

func (d docKeygenRequest) isCorrect() bool {
	return d.Credentials.IsCorrect() && d.ServerPubKey != ""
}

type docEncryptRequest struct {
	EncryptedKey string `json:"encryptedKey"`
	// Hex encoded document to encrypt
	HexSecret string `json:"hexSecret"`
}

func (d docEncryptRequest) isCorrect() bool {
	return d.EncryptedKey != "" && d.HexSecret != ""
}

type shadowDecryptRequest struct {
	parity.Credentials
	Secret       string   `json:"secret"`
	CommonPoint  string   `json:"commonPoint"`
	Shadows      []string `json:"shadows"`
	EncryptedDoc string   `json:"encryptedDoc"`
}

func (s shadowDecryptRequest) isCorrect() bool {
	return s.Credentials.IsCorrect() || s.Secret != "" ||
		s.CommonPoint != "" || s.EncryptedDoc != "" || len(s.Shadows) != 0
}

// -- Responses --

// SignRawHashResponse returns the successfully signed hash
type SignRawHashResponse struct {
	SignedHash string `json:"signedHash"`
}

// DocEncryptResponse returns the susccessfully encrypted doc
type DocEncryptResponse struct {
	EncryptedDoc string `json:"encryptedDoc"`
}

// ShadowDecryptResponse returns the successfully decrypted doc
type ShadowDecryptResponse struct {
	PlainDoc string `json:"plainDoc"`
}

// ----- FUNCTIONS ------

func signRandomHashHandler(c *gin.Context) {
	var cred parity.Credentials
	c.BindJSON(&cred)
	if !cred.IsCorrect() {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
	}
	doc, e := core.SignRandomHash(cred.Address, cred.Password)
	if e != nil {
		log.Println("SIGNRANDOMHASH FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, SignRawHashResponse{doc})
}

func signHashHandler(c *gin.Context) {
	var sr signRawHashRequest
	c.BindJSON(&sr)
	if sr.Address == "" || sr.Password == "" || sr.DocumentID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
	}
	doc, e := core.SignRawHash(sr.Address, sr.Password, sr.DocumentID)
	if e != nil {
		log.Println("SIGNHASH FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, SignRawHashResponse{doc})
}

func docKeyGenHandler(c *gin.Context) {

	var dr docKeygenRequest
	c.BindJSON(&dr)
	if dr.Address == "" || dr.Password == "" || dr.ServerPubKey == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
		return
	}
	docKeys, e := core.GenerateDocumentKey(dr.Address, dr.Password, dr.ServerPubKey)
	if e != nil {
		log.Println("DOCKEYGEN FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	log.Println("DOCKEYGEN SUCCESS ->", docKeys)

	c.JSON(http.StatusOK, docKeys)
}

func docEncryptHandler(c *gin.Context) {
	var dr docEncryptRequest
	c.BindJSON(&dr)
	if dr.EncryptedKey == "" || dr.HexSecret == "" {
		c.JSON(http.StatusBadRequest, "JSON parsing failed")
		return
	}

	encDoc, e := core.EncryptDocument("", "", dr.EncryptedKey, dr.HexSecret)
	if e != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, DocEncryptResponse{encDoc})
}

func shadowDecryptHandler(c *gin.Context) {
	var sr shadowDecryptRequest
	c.BindJSON(&sr)
	if !sr.isCorrect() {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing error"})
		return
	}

	hexDoc, e := core.DecryptDoc(sr.Address, sr.Password, sr.Secret, sr.CommonPoint, sr.Shadows, sr.EncryptedDoc)
	if e != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, ShadowDecryptResponse{hexDoc})
}
