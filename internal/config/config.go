package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var default_editor string = "vim"

type WikiConfig struct {
	Id               string   `json:"id"`
	Address          string   `json:"address"`
	User             string   `json:"user"`
	Password         string   `json:"password"`
	ImportExtensions []string `json:"importExtensions"`
}

type ConfigRoot struct {
	Wikis   []WikiConfig `json:"wikis"`
	Default string       `json:"default"`
	Editor  string       `json:"editor"`
}

type UserSettings struct {
	Editor string
}

var DefaultUserSettings *UserSettings = &UserSettings{
	Editor: default_editor,
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

func Get() (*WikiConfig, *ConfigRoot, error) {
	configRoot, err := getConfig()

	if err != nil {
		return &WikiConfig{}, &ConfigRoot{}, err
	}

	config, err := resolveDefaultConfig(configRoot)
	if err != nil {
		return &WikiConfig{}, &ConfigRoot{}, err
	}

	return config, configRoot, nil
}

func GetUserSettings(configRoot *ConfigRoot) *UserSettings {
	return &UserSettings{resolveUserEditor(configRoot)}
}

func resolveUserEditor(configRoot *ConfigRoot) string {
	editor := os.Getenv("EDITOR")

	if editor != "" {
		return editor
	}

	if configRoot.Editor != "" {
		return configRoot.Editor
	}

	return default_editor
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

func resolveDefaultConfig(configRoot *ConfigRoot) (*WikiConfig, error) {
	for _, config := range configRoot.Wikis {
		if configRoot.Default == config.Id {
			return &config, nil
		}
	}

	if len(configRoot.Wikis) == 0 {
		return &WikiConfig{}, errors.New("No configs found.")
	}

	return &configRoot.Wikis[0], nil
}

func getConfig() (*ConfigRoot, error) {
	configFilePath, configFileExists, err := GetConfigFilePath()
	if err != nil {
		return &ConfigRoot{}, err
	}

	if !configFileExists {
		return &ConfigRoot{}, ErrConfigDoesNotExist
	}

	return GetConfigFromPath(configFilePath)
}

func GetConfigFromPath(configFilePath string) (*ConfigRoot, error) {
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

func ImportExtensionsMedia(config *WikiConfig) []string {
	return append(imageExtensions, config.ImportExtensions...)
}
