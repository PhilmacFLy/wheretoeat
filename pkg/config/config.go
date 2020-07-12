package config

import (
	"encoding/json"
	"errors"
	"os"
)

//Config is the struct to save and load the config file
type Config struct {
	GoogleAPIKey string `json:"googleapikey"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Weight       Weight `json:"weight"`
}

//Weight is the struct to save the weights of the criteria
type Weight struct {
	Rating    float64 `json:"rating"`
	LastVisit float64 `json:"lastvisit"`
	DayCount  float64 `jsin:"daycount"`
}

//LoadConfig accepts a filepath and tries to load a config file from there
func LoadConfig(filepath string) (Config, error) {
	var res Config
	file, err := os.Open(filepath)
	if err != nil {
		return res, errors.New("Error opening file: " + err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		return res, errors.New("Error decoding file: " + err.Error())
	}
	return res, nil
}
