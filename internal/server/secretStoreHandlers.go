package server

import (
	"net/http"

	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/parity"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/core"
	"github.com/gin-gonic/gin"
)

// ----- STRUCTURES -----
// -- Requests --

type keyStoreRequest struct {
	parity.IDandSigned
	CommonPoint    string `json:"commonPoint"`
	EncryptedPoint string `json:"encryptedPoint"`
}

func (k keyStoreRequest) isCorrect() bool {
	return k.IDandSigned.IsCorrect() && k.CommonPoint != "" && k.EncryptedPoint != ""
}

// -- Responses --

// KeygenResponse returns the public portion of the Server Key
type KeygenResponse struct {
	Pubkey string `json:"serverPubkey"`
}

// ----- FUNCTIONS -----

func keygenHandler(c *gin.Context) {
	var is parity.IDandSigned
	c.BindJSON(&is)
	if !is.IsCorrect() {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
		return
	}

	pub, e := core.GetServerKey(is.DocID, is.SignedDocID)
	if e != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}
	c.JSON(http.StatusOK, KeygenResponse{pub})
}

func keyStoreHandler(c *gin.Context) {
	var kr keyStoreRequest
	c.BindJSON(&kr)
	if !kr.isCorrect() {
		c.JSON(http.StatusBadRequest, ErrorResponse{"JSON parsing failed"})
		return
	}

	e := core.StoreKey(kr.DocID, kr.SignedDocID, kr.CommonPoint, kr.EncryptedPoint)
	if e != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{e.Error()})
		return
	}

	c.String(http.StatusOK, "")
}
