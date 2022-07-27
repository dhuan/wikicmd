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

	testutils.RunWikiCmd(testState, "edit some_page", testutils.SetFakeVimToAddContent(" This is the new content."))

	testutils.MockAssert(
		t,
		&mock.AssertConfig{
			Route: "api.php",
			Nth:   5,
			Assert: &mock.AssertOptions{
				Type:  mock.AssertType_MethodMatch,
				Value: "post",
				And: &mock.AssertOptions{
					Type: mock.AssertType_FormMatch,
					KeyValues: map[string]interface{}{
						"action": "edit",
						"title":  "some_page",
						"text":   "This is a wiki page. This is the new content.",
					},
				},
			},
		},
	)
}
