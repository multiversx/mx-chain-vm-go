package main

import (
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	args := []string{"bin", "create-account", "--address=erdfoo", "--balance=100000", "--nonce=42"}
	app := initializeCLI(&arwendebug.DebugFacade{})
	err := app.Run(args)
	require.Nil(t, err)
}
