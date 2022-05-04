package wikicmd_test

import (
	"fmt"
	"testing"

	"github.com/dhuan/wikicmd/pkg/testutils"
)

func TestEditing(t *testing.T) {
	testState := testutils.StartupTest()
	commandResult, err := testutils.RunWikiCmd(testState, "edit")
	if err != nil {
		panic(err)
	}

	fmt.Println("111111111111111111")
	fmt.Println(commandResult)
}
