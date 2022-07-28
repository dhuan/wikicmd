package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/dhuan/mock/pkg/mock"
	"github.com/stretchr/testify/assert"
)

type TestState struct {
	WikicmdBinaryPath string
	WikicmdConfigPath string
	MockBinaryPath    string
}

type ConfigField int

const (
	Config_field_default ConfigField = iota
)

func RunWikiCmd(state *TestState, command string, env map[string]string) (string, error) {
	commandParameters := toCommandParameters(command)

	cmd := exec.Command(state.WikicmdBinaryPath, commandParameters...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("WIKICMD_CONFIG=%s", state.WikicmdConfigPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("EDITOR=%s/tests/bin/fakevim", pwd()))

	for key, _ := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, env[key]))
	}

	stdinBuffer := bytes.Buffer{}
	stdinBuffer.Write([]byte(fmt.Sprintln("no")))
	cmd.Stdin = &stdinBuffer

	result, err := cmd.CombinedOutput()
	if err != nil {
		return string(result), err
	}

	return trimEmptyLines(string(result)), nil
}

func trimEmptyLines(text string) string {
	lines := strings.Split(text, "\n")

	if strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	if strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
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
		fmt.Sprintf("%s/tests/bin/mock", pwd()),
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

type KillMockFunc func()

func parseCommandVars(command *string) {
	vars := map[string]string{
		"TEST_DATA_PATH": fmt.Sprintf("%s/tests/e2e/data", pwd()),
		"TEST_E2E_PORT":  "4000",
	}

	for key, value := range vars {
		*command = strings.Replace(
			*command,
			fmt.Sprintf("{{%s}}", key),
			value,
			-1,
		)
	}
}

func SetFakeVimToAddContent(content string) map[string]string {
	return map[string]string{
		"FAKEVIM_MODE":    "append",
		"FAKEVIM_CONTENT": content,
	}
}

func RunMockBg(state *TestState) KillMockFunc {
	command := fmt.Sprintf("serve -c %s/tests/e2e/mock/config/config.json -p 4000", pwd())
	parseCommandVars(&command)
	commandParameters := toCommandParameters(command)

	cmd := exec.Command(state.MockBinaryPath, commandParameters...)
	buf := &bytes.Buffer{}
	cmd.Stdout = buf
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	serverIsReady := waitForOutputInCommand("Mock server is listening on port 4000.", 4, buf)
	if !serverIsReady {
		panic("Something went wrong while waiting for mock to start up.")
	}

	return func() {
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	}
}

func waitForOutputInCommand(expectedOutput string, attempts int, buffer *bytes.Buffer) bool {
	for attempts > 0 {
		if strings.Contains(buffer.String(), expectedOutput) {
			return true
		}

		time.Sleep(500 * time.Millisecond)

		attempts--
	}

	return false
}

func MockAssert(t *testing.T, assertConfig *mock.AssertConfig) {
	mockConfig := mock.Init("localhost:4000")
	validationErrors, err := mock.Assert(mockConfig, assertConfig)
	if err != nil {
		panic(err)
	}

	if len(validationErrors) > 0 {
		fmt.Printf("Mock assertion failed!\n\n")

		for i, _ := range validationErrors {
			fmt.Printf("%+v\n\n", validationErrors[i])
		}

		t.Fail()
	}
}

func AssertLine(t *testing.T, lineNumber int, fullText, expectedText string) {
	lines := strings.Split(fullText, "\n")

	if lineNumber < 0 {
		lineNumber = len(lines) + lineNumber
	}

	assert.Equal(
		t,
		lines[lineNumber],
		expectedText,
	)
}
