package wikicmd_test

import (
	"testing"

	"github.com/dhuan/mock/pkg/mock"
	"github.com/dhuan/wikicmd/pkg/testutils"
)

func TestEdit(t *testing.T) {
	testState := testutils.StartupTest()
	killMock := testutils.RunMockBg(testState)
	defer killMock()

	testutils.RunWikiCmd(testState, "edit foobar", testutils.SetFakeVimToAddContent(" More content to this page."))

	testutils.MockAssert(
		t,
		&mock.AssertConfig{
			Route: "api.php",
			Assert: &mock.AssertOptions{
				Type:  mock.AssertType_MethodMatch,
				Value: "get",
			},
		},
	)
}
