package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func Edit(content string) (string, error) {
	if _, err := ioutil.TempFile("/var/tmp", "huahua"); err != nil {
		fmt.Println(err)
		panic("lascou")
	}

	if err := os.WriteFile("/var/tmp/huahua", []byte(content), 0644); err != nil {
		fmt.Println(err)
		panic("lascou")
	}

	cmd := exec.Command("vim", "/var/tmp/huahua")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic("eita")
	}

	fileContent, err := os.ReadFile("/var/tmp/huahua")

	if err != nil {
		panic("eita")
	}

	return string(fileContent), nil
}
