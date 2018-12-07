package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Run sets up the server on the given port
func Run(port int) {

	p := ":" + strconv.Itoa(port)

	router := gin.Default()
	router.GET("/insertRandomData", insertRandomDataHandler)
	router.GET("/signRandomHash", signRandomHashHandler)
	router.GET("/docAndKeygen", serverDocKeygenHandler)
	router.GET("/keygen", keygenHandler)
	router.POST("/decryptDataFromID", decryptDataFromIDHandler)

	router.Run(p)
}
