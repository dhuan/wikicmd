package wikicmd_test_unit

import (
	"testing"

	"github.com/dhuan/wikicmd/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestReplaceRegex(t *testing.T) {
	assert.Equal(
		t,
		"file",
		utils.ReplaceRegex(
			"file.txt",
			[]string{`\.txt$`},
			"",
		),
	)

	assert.Equal(
		t,
		"file",
		utils.ReplaceRegex(
			"file.wikitext",
			[]string{`\.go$`, `\.wikitext$`, `\.txt$`},
			"",
		),
	)
}

func TestFilePathToPageName(t *testing.T) {
	assert.Equal(
		t,
		"Some Page",
		utils.FilePathToPageName(
			[]string{"wikitext", "txt"},
			"Some Page.wikitext",
		),
	)

	assert.Equal(
		t,
		"Some Page",
		utils.FilePathToPageName(
			[]string{"wikitext", "txt"},
			"Some Page",
		),
	)
}
