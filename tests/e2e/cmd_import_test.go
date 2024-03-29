package wikicmd_test

import (
	"testing"

	"github.com/dhuan/mock/pkg/mock"
	"github.com/dhuan/wikicmd/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	testState := testutils.StartupTest()
	killMock := testutils.RunMockBg(testState)
	defer killMock()

	output, _ := testutils.RunWikiCmd(testState, "import pages_to_be_imported/page_one.txt", testutils.SetFakeVimToAddContent(" This is the new content."))

	assert.Equal(
		t,
		output,
		`Importing pages_to_be_imported/page_one.txt
1 item(s) have been imported.
Done!`,
	)

	testutils.MockAssert(
		t,
		&mock.AssertOptions{
			Route: "api.php",
			Nth:   4,
			Condition: &mock.Condition{
				Type:  mock.ConditionType_MethodMatch,
				Value: "post",
				And: &mock.Condition{
					Type: mock.ConditionType_FormMatch,
					KeyValues: map[string]interface{}{
						"action": "edit",
						"title":  "page_one",
						"text":   "This is page one. It will be imported during E2E tests.\n",
					},
				},
			},
		},
	)
}
