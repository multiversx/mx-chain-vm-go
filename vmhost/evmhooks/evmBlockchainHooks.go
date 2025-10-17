package evmhooks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/multiversx/mx-chain-core-go/core"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"math/big"
)

func (context *EVMHooksImpl) ChainID() *big.Int {
	return new(big.Int).SetBytes(context.GetBlockchainContext().ChainID())
}

func (context *EVMHooksImpl) Random() *common.Hash {
	random := common.BytesToHash(context.GetBlockchainContext().CurrentRandomSeed())
	return &random
}

func (context *EVMHooksImpl) GetHash(number uint64) common.Hash {
	return common.BytesToHash(context.GetBlockchainContext().BlockHash(number))
}

func (context *EVMHooksImpl) BlockNumber() *big.Int {
	return new(big.Int).SetUint64(context.GetBlockchainContext().CurrentNonce())
}

func (context *EVMHooksImpl) Time() uint64 {
	return context.GetBlockchainContext().CurrentTimeStamp()
}

func (context *EVMHooksImpl) GetBalanceForAddress(address []byte) *uint256.Int {
	return uint256.MustFromBig(context.GetBlockchainContext().GetBalanceBigInt(address))
}

func (context *EVMHooksImpl) GetSelfBalance() *uint256.Int {
	return context.GetBalanceForAddress(context.ContractMvxAddress())
}

func (context *EVMHooksImpl) GetBalance(address common.Address) *uint256.Int {
	return context.GetBalanceForAddress(context.toMVXAddress(address))
}

func (context *EVMHooksImpl) GetCodeHash(address common.Address) common.Hash {
	codeHash := context.GetBlockchainContext().GetCodeHash(context.toMVXAddress(address))
	return common.BytesToHash(codeHash)
}

func (context *EVMHooksImpl) GetCode(address common.Address) []byte {
	code, _ := context.GetBlockchainContext().GetCode(context.toMVXAddress(address))
	return code
}

func (context *EVMHooksImpl) GetCodeSize(address common.Address) int {
	return len(context.GetCode(address))
}

func (context *EVMHooksImpl) SaveAliasAddress() error {
	aliasAddress, err := context.requestEthereumContractAddress()
	if err != nil {
		return err
	}

	saveRequest := &vmcommon.AliasSaveRequest{
		AliasAddress:      aliasAddress.Bytes(),
		AliasIdentifier:   core.ETHAddressIdentifier,
		MultiversXAddress: context.ContractMvxAddress(),
	}
	return context.GetBlockchainContext().SaveAliasAddress(saveRequest)
}

func (context *EVMHooksImpl) requestEthereumContractAddress() (common.Address, error) {
	aliasAddress := context.ContractAliasAddress()
	if aliasAddress != (common.Address{}) {
		return aliasAddress, nil
	}

	return context.createEthereumContractAddress(context.CallerMvxAddress())
}

func (context *EVMHooksImpl) createEthereumContractAddress(creatorAddress []byte) (common.Address, error) {
	nonce, err := context.GetBlockchainContext().GetNonceForNewAddress(creatorAddress)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.CreateAddress(context.toEVMAddress(creatorAddress), nonce), nil
}
