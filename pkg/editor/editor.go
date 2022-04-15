package editor

import (
	"os"
	"os/exec"
)

var default_editor = "vim"

func Edit(content string) (string, bool, error) {
	fileName, err := mktemp()
	if err != nil {
		return "", false, err
	}

	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return "", false, err
	}

	fileInfoA, err := os.Stat(fileName)
	if err != nil {
		return "", false, err
	}

	cmd := exec.Command(getUserEditorCommand(), fileName)
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

	fileInfoB, err := os.Stat(fileName)
	if err != nil {
		return "", false, err
	}

	changed := fileInfoA.ModTime().String() != fileInfoB.ModTime().String()

	return string(fileContent), changed, nil
}

func EditFile(filePath string) error {
	cmd := exec.Command(getUserEditorCommand(), filePath)
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

func getUserEditorCommand() string {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		return default_editor
	}

	return editor
}
