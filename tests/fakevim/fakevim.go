package main

import (
	"fmt"
	"os"
)

type OperationMode int

const (
	OperationMode_None OperationMode = iota
	OperationMode_Append
	OperationMode_Overwrite
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No file was given! Check the manual.")

		os.Exit(1)
	}

	filePath := os.Args[1]
	operationMode := resolveOperationMode(os.Getenv("FAKEVIM_MODE"))

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	if operationMode == OperationMode_Append {
		err = writeFile(filePath, contentAppended(fileContent, os.Getenv("FAKEVIM_CONTENT")))
		if err != nil {
			panic(err)
		}
	}

	if operationMode == OperationMode_Overwrite {
		err = writeFile(filePath, []byte(os.Getenv("FAKEVIM_CONTENT")))
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("fakevim:%s\n", filePath)
}

func writeFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}

func contentAppended(contentBase []byte, contentMore string) []byte {
	return []byte(string(contentBase) + contentMore)
}

func resolveOperationMode(modeEncoded string) OperationMode {
	if modeEncoded == "" {
		return OperationMode_None
	}

	if modeEncoded == "append" {
		return OperationMode_Append
	}

	if modeEncoded == "overwrite" {
		return OperationMode_Overwrite
	}

	panic(fmt.Sprintf("Failed to resolve Operation Mode: %s", modeEncoded))
}
