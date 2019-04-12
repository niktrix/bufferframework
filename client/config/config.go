package config

import (
	"encoding/json"
	"os"

	"github.com/google/logger"
)

//Configuration Holds config struct
type Configuration struct {
	Server string  `json:"server"`
	Data   []int32 `json:"data"`
}

//Config Parse config.json and return config struct
func Init() (configuration *Configuration) {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		logger.Fatal("Error While opening config file", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		logger.Error(err)
	}
	return configuration
}
