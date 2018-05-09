package main

import (
	"os/user"
	"os"
	"encoding/json"
	"io/ioutil"
)


var config Config
var configDir string
var configPath string



type Config struct {
	Store string
}


func loadConfig(){

	usr, err := user.Current()
	ensure(err)

	configDir = usr.HomeDir + "/.timr"
	configPath = configDir + "/config.json"

	err = os.MkdirAll(configDir,os.ModePerm)
	ensure(err)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = Config{Store:"file"}
		saveConfig()
	}

	b, err := ioutil.ReadFile(configPath)
	ensure(err)

	json.Unmarshal(b, &config)
}


func saveConfig(){
	b, err := json.Marshal(config)
	ensure(err)

	err = ioutil.WriteFile(configPath, b, 0644)
	ensure(err)
}

