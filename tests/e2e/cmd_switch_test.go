package wikicmd_test

import (
	"testing"

	"github.com/dhuan/wikicmd/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestSwitchingToUnexistingWiki(t *testing.T) {
	testState := testutils.StartupTest()
	commandResult, _ := testutils.RunWikiCmd(testState, "switch unexisting_wiki")

	assert.Equal(
		t,
		commandResult,
		`No wiki exists with the given ID: unexisting_wiki

The available Wikis you can switch to are: my_wiki,another_wiki
`,
	)
}

func TestSwitchingToAnotherWiki(t *testing.T) {
	testState := testutils.StartupTest()
	commandResult, _ := testutils.RunWikiCmd(testState, "switch another_wiki")

	assert.Equal(
		t,
		commandResult,
		`Done!
`,
	)

	testutils.AssertConfig(t, testState, testutils.Config_field_default, "another_wiki")
}
