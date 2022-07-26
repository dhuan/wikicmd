package wikicmd_test

import (
	"fmt"
	"testing"

	"github.com/dhuan/wikicmd/pkg/testutils"
)

func TestEdit(t *testing.T) {
	testState := testutils.StartupTest()
	killMock := testutils.RunMockBg(testState)
	defer killMock()

	commandResult, _ := testutils.RunWikiCmd(testState, "edit foobar", testutils.SetFakeVimToAddContent(" More content to this page."))

	fmt.Println("!!!!!!!!!!!!!!")
	fmt.Println(commandResult)
}
