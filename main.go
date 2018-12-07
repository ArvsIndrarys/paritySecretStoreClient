package main

import (
	"log"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/config"
	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/server"
)

const (
	port           = 3333
	configFilePath = "./config.yml"
)

func main() {

	e := config.InitConfig(configFilePath)
	if e != nil {
		log.Println("init:", e)
		return
	}

	server.Run(port)
}
