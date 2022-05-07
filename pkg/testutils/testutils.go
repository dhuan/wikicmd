package testutils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestState struct {
	WikicmdBinaryPath string
	WikicmdConfigPath string
}

type ConfigField int

const (
	Config_field_default ConfigField = iota
)

func RunWikiCmd(state *TestState, command string) (string, error) {
	commandParameters := toCommandParameters(command)

	cmd := exec.Command(state.WikicmdBinaryPath, commandParameters...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("WIKICMD_CONFIG=%s", state.WikicmdConfigPath))
	result, err := cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}

	return string(result), nil
}

func StartupTest() *TestState {
	configFilePath, err := mktemp()
	if err != nil {
		panic(err)
	}

	configContent, err := os.ReadFile(fmt.Sprintf("%s/tests/e2e/wikicmd_config.json", pwd()))
	if err != nil {
		panic(err)
	}

	os.WriteFile(configFilePath, configContent, 0644)

	return &TestState{
		fmt.Sprintf("%s/bin/wikicmd", pwd()),
		configFilePath,
	}
}

func AssertConfig(t *testing.T, state *TestState, field ConfigField, expectedValue string) {
	var jsonParsed map[string]interface{}

	configContent, err := os.ReadFile(state.WikicmdConfigPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configContent, &jsonParsed)
	if err != nil {
		panic(err)
	}

	value := ""
	if field == Config_field_default {
		value = jsonParsed["default"].(string)
	}

	assert.Equal(
		t,
		value,
		expectedValue,
	)
}

func mktemp() (string, error) {
	result, err := exec.Command("mktemp").Output()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func pwd() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s/../..", wd)
}

func toCommandParameters(command string) []string {
	splitResult := strings.Split(command, " ")

	return splitResult
}
