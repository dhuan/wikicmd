package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type TestState struct {
	WikicmdBinaryPath string
	WikicmdConfigPath string
}

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
	return &TestState{
		fmt.Sprintf("%s/bin/wikicmd", pwd()),
		fmt.Sprintf("%s/tests/e2e/wikicmd_config.json", pwd()),
	}
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
