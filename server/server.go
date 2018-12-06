package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func server() {

	router := gin.Default()

	router.GET("/insertRandomData", insertRandomDataHandler)
	router.GET("/signRandomHash", signRandomHashHandler)
	router.GET("/docAndKeygen", serverDocKeygenHandler)
	router.GET("/keygen", keygenHandler)
	router.POST("/decryptDataFromID", decryptDataFromIDHandler)
	router.Run(port)
}

func insertRandomDataHandler(c *gin.Context) {

	docID, e := insertRandomDataInSecretStore()
	if e != nil {
		log.Println("INSERTION FAILURE:", e)
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	log.Println("INSERTION SUCCESS ->", docID)

	c.JSON(http.StatusOK, PublishResponse{docID})
}

func signRandomHashHandler(c *gin.Context) {
	doc, e := signRandomHash()
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	c.JSON(http.StatusOK, PublishResponse{doc})
}

func keygenHandler(c *gin.Context) {
	k, e := genRandomKey()
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
	}
	c.JSON(http.StatusOK, PublishResponse{k})
}

func serverDocKeygenHandler(c *gin.Context) {
	k, e := serverDocKeygen()
	if e != nil {
		c.JSON(http.StatusInternalServerError, e)
	}
	c.JSON(http.StatusOK, PublishResponse{k})
}

func decryptDataFromIDHandler(c *gin.Context) {

	var req DecryptRequestWithID
	c.BindJSON(&req)
	plainData, e := decryptViaStore(req.DocumentID)
	if e != nil {
		log.Println("DECRYPTION FAILURE:", e)
		c.JSON(http.StatusInternalServerError, e)
		return
	}
	log.Println("DECRYPTION SUCCESS ->", req.DocumentID)
	c.JSON(http.StatusOK, DecryptResponse{plainData})
}
