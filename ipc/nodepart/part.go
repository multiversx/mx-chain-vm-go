package nodepart

import (
	"fmt"
	"os"
	"time"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// NodePart is the endpoint that implements the message loop on Node's side
type NodePart struct {
	Logger     logger.Logger
	Messenger  *NodeMessenger
	blockchain vmcommon.BlockchainHook
	Repliers   []common.MessageReplier
	config     Config
}

// NewNodePart creates the Node part
func NewNodePart(
	mainLogger logger.Logger,
	dialogueLogger logger.Logger,
	input *os.File,
	output *os.File,
	blockchain vmcommon.BlockchainHook,
	config Config,
	marshalizer marshaling.Marshalizer,
) (*NodePart, error) {
	messenger := NewNodeMessenger(dialogueLogger, input, output, marshalizer)

	part := &NodePart{
		Logger:     mainLogger,
		Messenger:  messenger,
		blockchain: blockchain,
		config:     config,
	}

	part.Repliers = common.CreateReplySlots(part.noopReplier)
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

func (part *NodePart) noopReplier(message common.MessageHandler) common.MessageHandler {
	part.Logger.Error("noopReplier called")
	return common.CreateMessage(common.UndefinedRequestOrResponse)
}

// StartLoop runs the main loop
func (part *NodePart) StartLoop(request common.MessageHandler) (common.MessageHandler, error) {
	defer part.timeTrack(time.Now(), "[NODE] end of loop")

	part.Messenger.SendContractRequest(request)
	response, err := part.doLoop()
	if err != nil {
		part.Logger.Error("[NODE]: end of loop", "err", err)
	}

	part.Messenger.ResetDialogue()
	return response, err
}

// doLoop ends when processing the transaction ends or in the case of a critical failure
// Critical failure = Arwen timeouts or crashes
// The error result is set only in case of critical failure
func (part *NodePart) doLoop() (common.MessageHandler, error) {
	remainingMilliseconds := part.config.MaxLoopTime

	for {
		message, duration, err := part.Messenger.ReceiveHookCallRequestOrContractResponse(remainingMilliseconds)
		if err != nil {
			return nil, err
		}

		remainingMilliseconds -= duration
		if remainingMilliseconds < 0 {
			return nil, common.ErrArwenTimeExpired
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
	defer part.timeTrack(time.Now(), fmt.Sprintf("replyToHookCallRequest %s", request.GetKindName()))

	replier := part.Repliers[request.GetKind()]
	hookResponse := replier(request)
	err := part.Messenger.SendHookCallResponse(hookResponse)
	return err
}

// SendStopSignal sends a stop signal to Arwen
// Should only be used for tests!
func (part *NodePart) SendStopSignal() error {
	request := common.NewMessageStop()

	err := part.Messenger.Send(request)
	if err != nil {
		return err
	}

	part.Logger.Info("Node: sent stop signal to Arwen.")
	return nil
}

func (part *NodePart) timeTrack(start time.Time, message string) {
	elapsed := time.Since(start)
	part.Logger.Trace(message, "duration", elapsed)
}
