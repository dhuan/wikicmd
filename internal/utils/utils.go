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
	fileExtension := strings.ToLower(GetFileExtension(filePath))

	return AnyEquals(extensionList, fileExtension)
}

func GetFileExtension(filePath string) string {
	return GetNthWord(filePath, ".", -1, "unknown")
}

func GetNthWord(str, divisor string, index int, fallback string) string {
	splitResult := strings.Split(str, divisor)

	computedIndex := index
	if index < 0 {
		computedIndex = len(splitResult) + index
	}

	if computedIndex < 0 {
		return fallback
	}

	if computedIndex > (len(splitResult) - 1) {
		return fallback
	}

	return splitResult[computedIndex]
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

func FilePathToPageName(allowedExtensions []string, filePath string) string {
	allowedExtensionsRegex := make([]string, len(allowedExtensions))
	fileName := filepath.Base(filePath)

	for i, extension := range allowedExtensions {
		allowedExtensionsRegex[i] = fmt.Sprintf(`\.%s$`, extension)
	}

	return ReplaceRegex(fileName, allowedExtensionsRegex, "")
}

func ReplaceRegex(subject string, find []string, replaceWith string) string {
	if len(find) == 0 {
		return subject
	}

	re := regexp.MustCompile(find[0])

	return ReplaceRegex(
		re.ReplaceAllString(subject, replaceWith),
		find[1:],
		replaceWith,
	)
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

func StringToBool(str string) bool {
	if strings.ToLower(str) == "true" {
		return true
	}

	return false
}

func BoolToString(value bool) string {
	if value {
		return "true"
	}

	return "false"
}

func RemoveNth[T interface{}](list []T, iToRemove int) []T {
	if iToRemove > (len(list) - 1) {
		return list
	}

	newList := make([]T, 0)

	for i, value := range list {
		if i == iToRemove {
			continue
		}

		newList = append(newList, value)
	}

	return newList
}

func TrimEmptyLines(text string) string {
	lines := strings.Split(text, "\n")

	if strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}

	if strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}
