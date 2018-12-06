package main

import "log"

const (
	alicepwd             = "alicepwd"
	bobpwd               = "bobpwd"
	aliceAddr            = "0xb733e6b08e76559d5e9e55d5c56f10949c5017fd"
	bobAddr              = "0x7b51be7a60c022fba408e6a82a9bf69e71a18528"
	treshold             = "2"
	port                 = ":3333"
	signRawHashMethod    = "secretstore_signRawHash"
	generateDocKeyMethod = "secretstore_generateDocumentKey"
	encryptDocMethod     = "secretstore_encrypt"
	decryptDocMethod     = "secretstore_shadowDecrypt"
	shadow               = "shadow"
	version              = "2.0"
	id                   = 1
)

func main() {

	e := initConnexion(configFilePath)
	if e != nil {
		log.Println("init:", e)
		return
	}
}
