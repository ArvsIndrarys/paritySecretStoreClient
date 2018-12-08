package parity

import (
	"fmt"
	"strings"
)

// Credentials authenticates a function caller
type Credentials struct {
	Address  string `json:"address"`
	Password string `json:"pwd"`
}

// IsCorrect checks a Credentials correctness
func (c Credentials) IsCorrect() bool {
	return len(c.Address) == 42 && c.Address[0:2] == "0x" && c.Password != ""
}

// IDandSigned are the corresponding server Key ID and its signature by a user
type IDandSigned struct {
	DocID       string `json:"docId"`
	SignedDocID string `json:"signedDocId"`
}

// IsCorrect checks a IDandSigned correctness
func (i IDandSigned) IsCorrect() bool {
	return i.DocID != "" && i.SignedDocID != ""
}

// EncryptionKey contains the parameters returned by Parity Secret Sharing
// to encrypt a secret
type EncryptionKey struct {
	CommonPoint    string `json:"common_point"`
	EncryptedKey   string `json:"encrypted_key"`
	EncryptedPoint string `json:"encrypted_point"`
}

func (e EncryptionKey) String() string {

	var b strings.Builder
	fmt.Fprintf(&b, `Encryption Key containing :
		\nCommon point: %s,\nEncrypted key: %s\nEncrypted point: %s`,
		e.CommonPoint, e.EncryptedKey, e.EncryptedPoint)
	return b.String()
}

// DecryptionKey contains the parameters returned by Parity Secret Sharing
// to decrypt a secret
type DecryptionKey struct {
	Secret      string   `json:"decrypted_secret"`
	CommonPoint string   `json:"common_point"`
	Shadows     []string `json:"decrypt_shadows"`
}

// GetShadowsString allows recovering the secret Shadows as a string
func (k DecryptionKey) GetShadowsString() string {

	l := len(k.Shadows)
	var b strings.Builder
	fmt.Fprintf(&b, "[")
	for i := 0; i < l-1; i++ {
		fmt.Fprintf(&b, "\"%s\",", k.Shadows[i])
	}
	fmt.Fprintf(&b, "\"%s\"]", k.Shadows[l-1])
	return b.String()
}
