package testutils

import (
	"os/exec"
	"strings"
)

type TestState struct {
	WikicmdBinaryPath string
}

func RunWikiCmd(state *TestState, command string) (string, error) {
	commandParameters := toCommandParameters(command)

	result, err := exec.Command(state.WikicmdBinaryPath, commandParameters...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func StartupTest() *TestState {
	return &TestState{
		"/home/dev/github.com/dhuan/wikicmd/bin/wikicmd",
	}
}

func toCommandParameters(command string) []string {
	splitResult := strings.Split(command, " ")

	return splitResult
}
