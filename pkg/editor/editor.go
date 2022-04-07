package editor

import (
	"os"
	"os/exec"
)

func Edit(content string) (string, error) {
	fileName, err := mktemp()
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return "", err
	}

	cmd := exec.Command("vim", fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}

func EditFile(filePath string) error {
	cmd := exec.Command("vim", filePath)
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

	return string(result), nil
}
