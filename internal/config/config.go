package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Id               string   `json:"id"`
	Address          string   `json:"address"`
	User             string   `json:"user"`
	Password         string   `json:"password"`
	ImportExtensions []string `json:"importExtensions"`
}

type ConfigRoot struct {
	Config  []Config `json:"config"`
	Default string   `json:"default"`
}

func GetConfigFilePath() (string, bool, error) {
	customConfigPath := os.Getenv("WIKICMD_CONFIG")
	hasCustomConfigPath := customConfigPath != ""
	if hasCustomConfigPath {
		return customConfigPath, fileExists(customConfigPath), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", false, err
	}

	configFileName := fmt.Sprintf("%s/.wikicmd.json", homeDir)

	if !fileExists(configFileName) {
		return configFileName, false, nil
	}

	return configFileName, true, nil
}

func Get() (*Config, error) {
	configRoot, err := getConfig()

	if err != nil {
		return &Config{}, err
	}

	return resolveDefaultConfig(configRoot)
}

func Set(configRoot *ConfigRoot) error {
	configFilePath, _, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	fileContent, err := json.MarshalIndent(configRoot, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, fileContent, 0644)
}

func resolveDefaultConfig(configRoot *ConfigRoot) (*Config, error) {
	for _, config := range configRoot.Config {
		if configRoot.Default == config.Id {
			return &config, nil
		}
	}

	if len(configRoot.Config) == 0 {
		return &Config{}, errors.New("No configs found.")
	}

	return &configRoot.Config[0], nil
}

func getConfig() (*ConfigRoot, error) {
	configFilePath, configFileExists, err := GetConfigFilePath()
	if err != nil {
		return &ConfigRoot{}, err
	}

	if !configFileExists {
		return &ConfigRoot{}, errors.New("Config file not found.")
	}

	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return &ConfigRoot{}, err
	}

	decodedJson := &ConfigRoot{}

	if err = json.Unmarshal(fileContent, decodedJson); err != nil {
		return &ConfigRoot{}, err
	}

	return decodedJson, nil
}

func GetAll() (*ConfigRoot, error) {
	return getConfig()
}

func fileExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)

	if err != nil || fileInfo.IsDir() {
		return false
	}

	return true
}

func ImportExtensionsPage() []string {
	return pageExtensions
}

func ImportExtensionsMedia(config *Config) []string {
	return append(imageExtensions, config.ImportExtensions...)
}
