package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/core"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the result sent back in case of failure
type ErrorResponse struct {
	Err string `json:"error"`
}

// PublishResponse is the result of a publication request
type PublishResponse struct {
	DocumentID string `json:"documentId"`
}

// DecryptRequestWithID is the JSON that should be sent when requesting decryption
// of a secret with a given ID
type DecryptRequestWithID struct {
	DocumentID string `json:"documentId"`
}

// DecryptResponse is the result of a decryption request
type DecryptResponse struct {
	DecryptedData string `json:"plainData"`
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

func signRandomHashHandler(c *gin.Context) {
	doc, e := core.SignRandomHash()
	if e != nil {
		log.Println("SIGNRANDOMHASH FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, PublishResponse{doc})
}

func keygenHandler(c *gin.Context) {
	k, e := core.GenRandomKey()
	if e != nil {
		log.Println("KEYGEN FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, PublishResponse{k})
}

func serverDocKeygenHandler(c *gin.Context) {
	k, e := core.ServerDocKeygen()
	if e != nil {
		log.Println("DOCANDKEYGEN FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, PublishResponse{k})
}

func decryptDataFromIDHandler(c *gin.Context) {
	empty := DecryptRequestWithID{}

	var req DecryptRequestWithID
	c.BindJSON(&req)
	if req == empty {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON Parsing failed"})
		return
	}
	fmt.Println(req)
	plainData, e := core.DecryptViaStore(req.DocumentID)
	if e != nil {
		log.Println("DECRYPTION FAILURE:", e)
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	log.Println("DECRYPTION SUCCESS ->", req.DocumentID)
	c.JSON(http.StatusOK, DecryptResponse{plainData})
}
