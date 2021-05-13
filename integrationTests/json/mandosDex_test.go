package vmjsonintegrationtest

import (
	"testing"
)

func TestDex_v0_1(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "dex/v0_1")
}
