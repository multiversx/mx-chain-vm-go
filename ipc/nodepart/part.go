package nodepart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// NodePart is the endpoint that implements the message loop on Node's side
type NodePart struct {
	Logger     logger.Logger
	Messenger  *NodeMessenger
	blockchain vmcommon.BlockchainHook
	Repliers   []common.MessageReplier
}

// NewNodePart creates the Node part
func NewNodePart(nodeLogger logger.Logger, input *os.File, output *os.File, blockchain vmcommon.BlockchainHook) (*NodePart, error) {
	messenger := NewNodeMessenger(nodeLogger, input, output)

	part := &NodePart{
		Logger:     nodeLogger,
		Messenger:  messenger,
		blockchain: blockchain,
	}

	part.Repliers = common.CreateReplySlots()
	part.Repliers[common.BlockchainAccountExistsRequest] = part.replyToBlockchainAccountExists
	part.Repliers[common.BlockchainNewAddressRequest] = part.replyToBlockchainNewAddress
	part.Repliers[common.BlockchainGetBalanceRequest] = part.replyToBlockchainGetBalance
	part.Repliers[common.BlockchainGetNonceRequest] = part.replyToBlockchainGetNonce
	part.Repliers[common.BlockchainGetStorageDataRequest] = part.replyToBlockchainGetStorageData
	part.Repliers[common.BlockchainIsCodeEmptyRequest] = part.replyToBlockchainIsCodeEmpty
	part.Repliers[common.BlockchainGetCodeRequest] = part.replyToBlockchainGetCode
	part.Repliers[common.BlockchainGetBlockhashRequest] = part.replyToBlockchainGetBlockhash
	part.Repliers[common.BlockchainLastNonceRequest] = part.replyToBlockchainLastNonce
	part.Repliers[common.BlockchainLastRoundRequest] = part.replyToBlockchainLastRound
	part.Repliers[common.BlockchainLastTimeStampRequest] = part.replyToBlockchainLastTimeStamp
	part.Repliers[common.BlockchainLastRandomSeedRequest] = part.replyToBlockchainLastRandomSeed
	part.Repliers[common.BlockchainLastEpochRequest] = part.replyToBlockchainLastEpoch
	part.Repliers[common.BlockchainGetStateRootHashRequest] = part.replyToBlockchainGetStateRootHash
	part.Repliers[common.BlockchainCurrentNonceRequest] = part.replyToBlockchainCurrentNonce
	part.Repliers[common.BlockchainCurrentRoundRequest] = part.replyToBlockchainCurrentRound
	part.Repliers[common.BlockchainCurrentTimeStampRequest] = part.replyToBlockchainCurrentTimeStamp
	part.Repliers[common.BlockchainCurrentRandomSeedRequest] = part.replyToBlockchainCurrentRandomSeed
	part.Repliers[common.BlockchainCurrentEpochRequest] = part.replyToBlockchainCurrentEpoch

	return part, nil
}

// StartLoop runs the main loop
func (part *NodePart) StartLoop(request common.MessageHandler) (common.MessageHandler, error) {
	part.Messenger.SendContractRequest(request)
	response, err := part.doLoop()

	part.Logger.Debug("[NODE]: end of loop", "err", err)
	part.Messenger.ResetDialogue()
	return response, err
}

// doLoop ends when processing the transaction ends or in the case of a critical failure
// Critical failure = Arwen timeouts or crashes
// The error result is set only in case of critical failure
func (part *NodePart) doLoop() (common.MessageHandler, error) {
	const MaxLoopTime = 1000
	remainingMilliseconds := MaxLoopTime

	for {
		message, duration, err := part.Messenger.ReceiveHookCallRequestOrContractResponse(remainingMilliseconds)
		remainingMilliseconds -= duration
		if err != nil {
			return nil, err
		}

		if common.IsHookCall(message) {
			err := part.replyToHookCallRequest(message)
			if err != nil {
				return nil, err
			}

			continue
		}

		if common.IsContractResponse(message) {
			return message, nil
		}
		if common.IsDiagnose(message) {
			return message, nil
		}

		return nil, common.ErrBadMessageFromArwen
	}
}

func (part *NodePart) replyToHookCallRequest(request common.MessageHandler) error {
	replier := part.Repliers[request.GetKind()]
	hookResponse := replier(request)
	err := part.Messenger.SendHookCallResponse(hookResponse)
	return err
}

// SendStopSignal sends a stop signal to Arwen
// Should only be used for tests!
func (part *NodePart) SendStopSignal() error {
	request := common.NewMessageStop()

	err := part.Messenger.SendContractRequest(request)
	if err != nil {
		return err
	}

	part.Logger.Info("Node: sent stop signal to Arwen.")
	return nil
}
