package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v1"
)

const (
	configFilePath = "./config.yml"
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

func initConnexion(path string) error {
	f, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}

	e = yaml.Unmarshal(f, &config)
	if e != nil {
		return e
	}

	baseSecretStoreURL = URL{BaseURL: config.SecretStoreServerURL,
		Port: config.SecretStoreServerPort}

	baseSecretStoreMethodsURL = URL{BaseURL: config.ClientURL,
		Port: config.ClientPort}
	return nil
}
