package nodepart

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// NodePart is
type NodePart struct {
	Messenger  *NodeMessenger
	blockchain vmcommon.BlockchainHook
	Repliers   []common.MessageReplier
}

// NewNodePart creates
func NewNodePart(input *os.File, output *os.File, blockchain vmcommon.BlockchainHook) (*NodePart, error) {
	messenger := NewNodeMessenger(input, output)

	part := &NodePart{
		Messenger:  messenger,
		blockchain: blockchain,
	}

	part.Repliers = common.CreateReplySlots()
	part.Repliers[common.BlockchainNewAddressRequest] = part.replyToBlockchainNewAddress
	part.Repliers[common.BlockchainGetNonceRequest] = part.replyToBlockchainGetNonce
	part.Repliers[common.BlockchainGetStorageDataRequest] = part.replyToBlockchainGetStorageData
	part.Repliers[common.BlockchainGetCodeRequest] = part.replyToBlockchainGetCode

	return part, nil
}

// StartLoop runs the main loop
func (part *NodePart) StartLoop(request common.MessageHandler) (common.MessageHandler, error) {
	part.Messenger.SendContractRequest(request)
	response, err := part.doLoop()

	common.LogDebug("[NODE]: end of loop, err=%v", err)
	part.Messenger.EndDialogue()
	return response, err
}

// doLoop ends when processing the transaction ends or in the case of a critical failure
// Critical failure = Arwen timeouts or crashes
// The error result is set only in case of critical failure
func (part *NodePart) doLoop() (common.MessageHandler, error) {
	for {
		// TODO: start with initial timeout, decrement with "time.Since".
		// TODO: Allow a total max of 1 second (accumulated wait).
		message, err := part.Messenger.ReceiveHookCallRequestOrContractResponse(1000)
		if err != nil {
			return nil, err
		}

		if common.IsHookCallRequest(message) {
			err := part.replyToHookCallRequest(message)
			if err != nil {
				return nil, err
			}

			continue
		}

		if common.IsContractResponse(message) {
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

	common.LogDebug("Node: sent stop signal to Arwen.")
	return nil
}

func (part *NodePart) replyToBlockchainNewAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainNewAddressRequest)
	address, err := part.blockchain.NewAddress(typedRequest.CreatorAddress, typedRequest.CreatorNonce, typedRequest.VMType)
	response := common.NewMessageBlockchainNewAddressResponse(err)
	response.Address = address
	return response
}

func (part *NodePart) replyToBlockchainGetNonce(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetNonceRequest)
	nonce, err := part.blockchain.GetNonce(typedRequest.Address)
	response := common.NewMessageBlockchainGetNonceResponse(err)
	response.Nonce = nonce
	return response
}

func (part *NodePart) replyToBlockchainGetStorageData(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetStorageDataRequest)
	data, err := part.blockchain.GetStorageData(typedRequest.Address, typedRequest.Index)
	response := common.NewMessageBlockchainGetStorageDataResponse(err)
	response.Data = data
	return response
}

func (part *NodePart) replyToBlockchainGetCode(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetCodeRequest)
	code, err := part.blockchain.GetCode(typedRequest.Address)
	response := common.NewMessageBlockchainGetCodeResponse(err)
	response.Code = code
	return response
}
