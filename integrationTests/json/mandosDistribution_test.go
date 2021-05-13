package vmjsonintegrationtest

import (
	"testing"
)

func TestDistribution_v0_1(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "distribution/v0_1")
}
