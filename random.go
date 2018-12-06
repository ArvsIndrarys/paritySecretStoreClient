package main

import (
	"encoding/hex"
	"math/rand"
	"time"
)

const (
	letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	docIDlength = 32
)

func randDocID() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, docIDlength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	id := hex.EncodeToString(b)

	//e := writeFile(string(b), inputFolder+id)
	//if e != nil {
	//	log.Println(e)
	//	return ""
	//}

	return id
}
