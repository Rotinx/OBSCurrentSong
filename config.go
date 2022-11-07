package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type configStructure struct {
	SavePath string `json:"savePath"`
}

var config = configStructure{
	SavePath: "./data",
}

func (c configStructure) loadConfig() {
	err := os.MkdirAll("./data", 0644)
	if err != nil {
		panic(err)
		return
	}
	file, err := os.Open("./data/config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println(err)
			return
		}
		panic(err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	configuration := configStructure{}
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	config = configuration
}

func (c configStructure) save() {
	s, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("./data/config.json", s, 0644)
	if err != nil {
		panic(err)
	}
}
