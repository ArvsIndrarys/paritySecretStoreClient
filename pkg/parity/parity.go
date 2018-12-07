package parity

import (
	"fmt"
	"strings"
)

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

	var b strings.Builder
	fmt.Fprintf(&b, "[")
	for _, v := range k.Shadows {
		fmt.Fprintf(&b, v)
	}
	fmt.Fprintf(&b, "[")

	return b.String()
}
