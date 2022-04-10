package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func ValidateFiles(filePaths []string, acceptedExtensions []string) map[string]error {
	validationErrors := make(map[string]error)

	for _, filePath := range filePaths {
		if !FileExists(filePath) {
			validationErrors[filePath] = ErrFileDoesNotExist

			continue
		}

		if !ExtensionMatches(acceptedExtensions, filePath) {
			validationErrors[filePath] = ErrExtensionNotAccepted

			continue
		}
	}

	return validationErrors
}

func FileExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)

	if err != nil || fileInfo.IsDir() {
		return false
	}

	return true
}

func FilePathToPageName(filePath string) string {
	fileName := filepath.Base(filePath)

	return strings.Replace(fileName, ".wikitext", "", 1)
}

func Wget(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func AnyEquals[T comparable](list []T, value T) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}
