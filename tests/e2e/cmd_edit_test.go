package wikicmd_test

import (
	"testing"

	"github.com/dhuan/mock/pkg/mock"
	"github.com/dhuan/wikicmd/pkg/testutils"
)

func TestEditPage(t *testing.T) {
	testState := testutils.StartupTest()
	killMock := testutils.RunMockBg(testState)
	defer killMock()

	output, _ := testutils.RunWikiCmd(testState, "edit some_page", testutils.SetFakeVimToAddContent(" This is the new content."))

	testutils.MockAssert(
		t,
		&mock.AssertOptions{
			Route: "api.php",
			Nth:   5,
			Condition: &mock.Condition{
				Type:  mock.ConditionType_MethodMatch,
				Value: "post",
				And: &mock.Condition{
					Type: mock.ConditionType_FormMatch,
					KeyValues: map[string]interface{}{
						"action": "edit",
						"title":  "some_page",
						"text":   "This is a wiki page. This is the new content.",
					},
				},
			},
		},
	)

	testutils.AssertLine(t, -1, output, "some_page edited successfully.")
}
