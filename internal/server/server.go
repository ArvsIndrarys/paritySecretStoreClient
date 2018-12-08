package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Run sets up the server on the given port
func Run(port int) {

	p := ":" + strconv.Itoa(port)

	router := gin.Default()

	// -- Secret Store functions --
	router.POST("/keygen", keygenHandler)
	router.POST("/keyStore", keyStoreHandler)
	// ShadowKeyRegen GET
	// KeyRegen GET
	// Schnorr POST
	// ECDSASign POST

	// -- Parity Node functions --
	router.POST("/signRandomHash", signRandomHashHandler)
	router.POST("/docKeyGen", docKeyGenHandler)
	router.POST("/docEncrypt", docEncryptHandler)
	router.POST("/shadowDecrypt", shadowDecryptHandler)

	// -- Combinations --

	// handles the signRawHash and Keygen functions
	router.GET("/toRandKeygen", toRandKeygenHandler)
	router.POST("/toKeygen", toKeygenHandler)

	// handles the signRawHash and DocAndKeygen functions
	router.GET("/randDocAndKeygen", randDocAndKeygenHandler) // POST

	// full insertion scenario
	router.GET("/insertRandomData", insertRandomDataHandler) // POST
	router.GET("/insertData/:docID", insertDataHandler)      // POST
	// full decryption scenario from docID
	router.POST("/decryptDataFromID", decryptDataFromIDHandler) // POST

	router.Run(p)
}
