package config

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	Elasticsearch struct {
		Hosts []string `json:"hosts"`
	} `json:"elasticsearch"`
	Server  string `json:"server"`
	QaIndex string `json:"qaIndex"`
}

var Config = &configuration{}

func init() {
	file, e := os.Open("configs/config.json")
	if e != nil {
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if e := decoder.Decode(Config); e != nil {
		log.Println(e)
	}
	log.Println(*Config)
}
