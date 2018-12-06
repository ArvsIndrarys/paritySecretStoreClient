package main

// Query is a base Parity query
type Query struct {
	JSONRPCVersion string   `json:"jsonrpc"`
	Method         string   `json:"method"`
	Params         []string `json:"params"`
	ID             int      `json:"id"`
}

// EncKeyQueryResult is a Parity response containing an encryption Key
type EncKeyQueryResult struct {
	JSONRPCVersion string        `json:"jsonrpc"`
	Result         EncryptionKey `json:"result"`
	ID             int           `json:"id"`
}

// EncryptionKey contains the parameters returned by Parity Secret Sharing
// to encrypt a secret
type EncryptionKey struct {
	CommonPoint    string `json:"common_point"`
	EncryptedKey   string `json:"encrypted_key"`
	EncryptedPoint string `json:"encrypted_point"`
}

func (e EncryptionKey) String() string {
	return buildString("Encryption Key containing :",
		"\nCommon point: ", e.CommonPoint,
		"\nEncrypted key:", e.EncryptedKey,
		"\nEncrypted point:", e.EncryptedPoint)
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

	s := buildString("[")
	for _, v := range k.Shadows {
		s = buildString(s, "\"", v, "\"", ",")
	}
	s = buildString(s[:len(s)-1], "]")

	return s
}
