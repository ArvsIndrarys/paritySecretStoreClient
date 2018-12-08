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

type toKeygenRequest struct {
	parity.Credentials
	DocID       string `json:"docId"`
	SignedDocID string `json:"signedDocId"`
}

func (k toKeygenRequest) isCorrect() bool {
	return k.Credentials.IsCorrect() || k.DocID != "" || k.SignedDocID != ""
}

type decryptRequestWithID struct {
	DocumentID string `json:"documentId"`
}

// ErrorResponse is the result sent back in case of failure
type ErrorResponse struct {
	Err string `json:"error"`
}

// PublishResponse is the result of a publication request
type PublishResponse struct {
	DocumentID string `json:"documentId"`
}

// DecryptResponse is the result of a decryption request
type DecryptResponse struct {
	DecryptedData string `json:"plainData"`
}

func toRandKeygenHandler(c *gin.Context) {
	var cred parity.Credentials
	c.BindJSON(&cred)
	if !cred.IsCorrect() {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
		return
	}

	k, e := core.GenRandomKey()
	if e != nil {
		log.Println("KEYGEN FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, KeygenResponse{k})
}

func toKeygenHandler(c *gin.Context) {

	var kr toKeygenRequest
	c.BindJSON(&kr)
	if !kr.isCorrect() {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
		return
	}

	if kr.DocID == "" || kr.SignedDocID == "" {
		c.JSON(http.StatusBadRequest, "JSON Parsing failed")
		return
	}
	k, e := core.GetServerKey(kr.DocID, kr.SignedDocID)
	if e != nil {
		log.Println("KEYGEN FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, KeygenResponse{k})
}

func insertRandomDataHandler(c *gin.Context) {

	docID, e := core.InsertRandomDataInSecretStore()
	if e != nil {
		log.Println("INSERTION FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	log.Println("INSERTION SUCCESS ->", docID)

	c.JSON(http.StatusOK, PublishResponse{docID})
}

func insertDataHandler(c *gin.Context) {

	docID := c.Param("docID")
	e := core.InsertDataInSecretStore(docID)
	if e != nil {
		log.Println("INSERTION FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	log.Println("INSERTION SUCCESS ->", docID)

	c.JSON(http.StatusOK, PublishResponse{docID})
}

func decryptDataFromIDHandler(c *gin.Context) {
	empty := decryptRequestWithID{}

	var req decryptRequestWithID
	c.BindJSON(&req)
	if req == empty {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON Parsing failed"})
		return
	}
	plainData, e := core.DecryptViaStore(req.DocumentID)
	if e != nil {
		log.Println("DECRYPTION FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	log.Println("DECRYPTION SUCCESS ->", req.DocumentID)
	c.JSON(http.StatusOK, DecryptResponse{plainData})
}

func randDocAndKeygenHandler(c *gin.Context) {
	k, e := core.RandDocAndKeygen()
	if e != nil {
		log.Println("DOCANDKEYGEN FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, PublishResponse{k})
}
