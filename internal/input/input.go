package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func TextPrompt(prompt, fallback string) string {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if text == "" || text == "\n" {
		return fallback
	}

	return strings.ReplaceAll(text, "\n", "")
}

func ResponsePrompt[T interface{}](
	prompt string,
	choices map[string]T,
	fallback T,
	onEmpty T,
) (T, error) {
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		return fallback, err
	}

	textParsed := strings.TrimSpace(text)
	message, ok := choices[textParsed]

	if textParsed == "" {
		return onEmpty, nil
	}

	if !ok {
		return fallback, nil
	}

	return message, nil
}
