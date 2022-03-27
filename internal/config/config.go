package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Id       string `json:"id"`
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ConfigRoot struct {
	Config []Config `json:"config"`
}

func Get() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	configFileName := fmt.Sprintf("%s/.wikicmd.json", homeDir)

	if !fileExists(configFileName) {
		return Config{}, errors.New("Config file not found.")
	}

	fileContent, err := os.ReadFile(configFileName)
	if err != nil {
		return Config{}, err
	}

	decodedJson := &ConfigRoot{}

	if err = json.Unmarshal(fileContent, decodedJson); err != nil {
		return Config{}, err
	}

	if len(decodedJson.Config) == 0 {
		return Config{}, errors.New("No configs found.")
	}

	return decodedJson.Config[0], nil
}

func fileExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)

	if err != nil || fileInfo.IsDir() {
		return false
	}

	return true
}
