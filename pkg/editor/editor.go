package editor

import (
	"os"
	"os/exec"
	"reflect"

	"github.com/dhuan/wikicmd/internal/utils"
)

func Edit(editorProgram, content string) (string, bool, error) {
	fileName, err := mktemp()
	if err != nil {
		return "", false, err
	}

	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return "", false, err
	}

	fileContentBeforeEdited, err := os.ReadFile(fileName)
	if err != nil {
		return "", false, err
	}

	cmd := exec.Command(editorProgram, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", false, err
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return "", false, err
	}

	changed := !(reflect.DeepEqual(fileContentBeforeEdited, fileContent))

	return string(fileContent), changed, nil
}

func EditFile(editorProgram, filePath string) error {
	cmd := exec.Command(editorProgram, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func mktemp() (string, error) {
	result, err := exec.Command("mktemp").Output()
	if err != nil {
		return "", err
	}

	return utils.TrimEmptyLines(string(result)), nil
}
