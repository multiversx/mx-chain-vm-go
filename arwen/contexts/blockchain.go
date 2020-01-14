package subcontexts

import (
	"fmt"
	"math/big"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

type Blockchain struct {
	host           arwen.VMContext
	blockChainHook vmcommon.BlockchainHook
}

func NewBlockchainSubcontext(
	host arwen.VMContext,
	blockChainHook vmcommon.BlockchainHook,
) (*Blockchain, error) {

	blockchain := &Blockchain{
		blockChainHook: blockChainHook,
		host:           host,
	}

	return blockchain, nil
}

func (blockchain *Blockchain) NewAddress(creatorAddress []byte) ([]byte, error) {
	nonce, err := blockchain.GetNonce(creatorAddress)
	if err != nil {
		return nil, err
	}

	if nonce > 0 {
		nonce -= 1
	}

	vmType := blockchain.host.Runtime().GetVMType()
  return blockchain.blockChainHook.NewAddress(creatorAddress, nonce, vmType)
}

func (blockchain *Blockchain) AccountExists(addr []byte) bool {
	exists, err := blockchain.blockChainHook.AccountExists(addr)
	if err != nil {
		fmt.Printf("Account exsits returned with error %s \n", err.Error())
	}
	return exists
}

func (blockchain *Blockchain) GetBalance(addr []byte) []byte {
	strAdr := string(addr)

	outputAccounts := blockchain.host.Output().GetOutputAccounts()
	if _, ok := outputAccounts[strAdr]; ok {
		balance := outputAccounts[strAdr].Balance
		return balance.Bytes()
	}

	balance, err := blockchain.blockChainHook.GetBalance(addr)
	if err != nil {
		fmt.Printf("GetBalance returned with error %s \n", err.Error())
		return big.NewInt(0).Bytes()
	}

	outputAccounts[strAdr] = &vmcommon.OutputAccount{
		Balance:      big.NewInt(0).Set(balance),
		BalanceDelta: big.NewInt(0),
		Address:      addr,
	}

	return balance.Bytes()
}

func (blockchain *Blockchain) GetNonce(addr []byte) (uint64, error) {
	strAdr := string(addr)
	outputAccounts := blockchain.host.Output().GetOutputAccounts()
	if _, ok := outputAccounts[strAdr]; ok {
		return outputAccounts[strAdr].Nonce, nil
	}

	nonce, err := blockchain.blockChainHook.GetNonce(addr)
	if err != nil {
		fmt.Printf("GetNonce returned with error %s \n", err.Error())
	}

	outputAccounts[strAdr] = &vmcommon.OutputAccount{BalanceDelta: big.NewInt(0), Address: addr, Nonce: nonce}
	return nonce, err
}

func (blockchain *Blockchain) IncreaseNonce(addr []byte) {
	nonce, _ := blockchain.GetNonce(addr)
	outputAccounts := blockchain.host.Output().GetOutputAccounts()
	outputAccounts[string(addr)].Nonce = nonce + 1
}


func (blockchain *Blockchain) GetCodeHash(addr []byte) ([]byte, error) {
	code, err := blockchain.blockChainHook.GetCode(addr)
	if err != nil {
		return nil, err
	}

	return blockchain.host.Crypto().Keccak256(code)
}

func (blockchain *Blockchain) GetCode(addr []byte) ([]byte, error) {
	return blockchain.blockChainHook.GetCode(addr)
}

func (blockchain *Blockchain) GetCodeSize(addr []byte) (int32, error) {
	code, err := blockchain.blockChainHook.GetCode(addr)
	if err != nil {
		return 0, err
	}

	result := int32(len(code))
	return result, nil
}

func (blockchain *Blockchain) BlockHash(number int64) []byte {
	if number < 0 {
		fmt.Printf("BlockHash nonce cannot be negative\n")
		return nil
	}
	block, err := blockchain.blockChainHook.GetBlockhash(uint64(number))

	if err != nil {
		fmt.Printf("GetBlockHash returned with error %s \n", err.Error())
		return nil
	}

	return block
}

func (blockchain *Blockchain) CurrentEpoch() uint32 {
  return blockchain.blockChainHook.CurrentEpoch()
}

func (blockchain *Blockchain) CurrentNonce() uint64 {
  return blockchain.blockChainHook.CurrentNonce()
}

func (blockchain *Blockchain) GetStateRootHash() []byte {
	return blockchain.blockChainHook.GetStateRootHash()
}

func (blockchain *Blockchain) LastTimeStamp() uint64 {
	return blockchain.blockChainHook.LastTimeStamp()
}

func (blockchain *Blockchain) LastNonce() uint64 {
	return blockchain.blockChainHook.LastNonce()
}

func (blockchain *Blockchain) LastRound() uint64 {
	return blockchain.blockChainHook.LastRound()
}

func (blockchain *Blockchain) LastEpoch() uint32 {
	return blockchain.blockChainHook.LastEpoch()
}

func (blockchain *Blockchain) CurrentRound() uint64 {
	return blockchain.blockChainHook.CurrentRound()
}

func (blockchain *Blockchain) CurrentTimeStamp() uint64 {
	return blockchain.blockChainHook.CurrentTimeStamp()
}

func (blockchain *Blockchain) LastRandomSeed() []byte {
	return blockchain.blockChainHook.LastRandomSeed()
}

func (blockchain *Blockchain) CurrentRandomSeed() []byte {
  return blockchain.blockChainHook.CurrentRandomSeed()
}
