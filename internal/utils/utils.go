package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func FormatPageNameInput(pageName string) string {
	return strings.ReplaceAll(pageName, " ", "_")
}

func ExtensionMatches(extensionList []string, filePath string) bool {
	extensionListRegex := make([]string, len(extensionList))

	for i, extension := range extensionList {
		extensionListRegex[i] = fmt.Sprintf(".%s$", extension)
	}

	return RegexTestAny(extensionListRegex, filePath)
}

func RegexTestAny(regexList []string, subject string) bool {
	for _, regex := range regexList {
		if RegexTest(regex, subject) {
			return true
		}
	}

	return false
}

func RegexTest(regex string, subject string) bool {
	match, err := regexp.MatchString(regex, subject)

	if err != nil {
		return false
	}

	return match
}

func MapValueSearch[TMapKey comparable, TMapValue comparable](
	subject map[TMapKey]TMapValue,
	valueToBeSearched TMapValue,
	fallback TMapKey,
) TMapKey {
	for key, value := range subject {
		if value == valueToBeSearched {
			return key
		}
	}

	return fallback
}
