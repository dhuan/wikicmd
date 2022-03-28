package utils

import "strings"

func FormatPageNameInput(pageName string) string {
	return strings.ReplaceAll(pageName, " ", "_")
}
