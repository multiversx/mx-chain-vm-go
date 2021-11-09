package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultisig_DeployFactorial(t *testing.T) {
	err := runSingleTestReturnError("multisig/mandos", "deployFactorial.scen.json")
	require.Nil(t, err)
}

func TestDelegation_ActivateOtherShard(t *testing.T) {
	err := runSingleTestReturnError("delegation/v0_2/activate", "activate_other_shard.scen.json")
	require.Nil(t, err)
}
