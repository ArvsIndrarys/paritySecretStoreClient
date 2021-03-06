package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v1"

	"github.com/ArvsIndrarys/paritySecretStoreClient/internal/core"
	"github.com/ArvsIndrarys/paritySecretStoreClient/pkg/net"
)

var (
	config Config
)

// Config is the server connexion
type Config struct {
	SecretStoreServerURL  string `yaml:"secretstore_url"`
	SecretStoreServerPort string `yaml:"secretstore_port"`
	ClientURL             string `yaml:"client_url"`
	ClientPort            string `yaml:"client_port"`
}

func InitConfig(path string) error {
	f, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}

	e = yaml.Unmarshal(f, &config)
	if e != nil {
		return e
	}

	core.BaseSecretStoreURL = net.URL{BaseURL: config.SecretStoreServerURL,
		Port: config.SecretStoreServerPort}

	core.BaseSecretStoreMethodsURL = net.URL{BaseURL: config.ClientURL,
		Port: config.ClientPort}
	return nil
}
