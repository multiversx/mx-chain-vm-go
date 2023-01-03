package testcommon

import (
	"sort"

	"github.com/ElrondNetwork/elrond-go-core/data/vm"
)

// CrossShardCall -
type CrossShardCall struct {
	CallerAddress []byte
	StartNode     *TestCallNode
	CallType      vm.CallType
	AsyncData     []byte
	Data          []byte
	ParentsPath   []*TestCallNode
}

// CrossShardCallsQueue -
type CrossShardCallsQueue struct {
	Data []*CrossShardCall
}

// NewCrossShardCallQueue -
func NewCrossShardCallQueue() *CrossShardCallsQueue {
	return &CrossShardCallsQueue{}
}

// Enqueue -
func (queue *CrossShardCallsQueue) Enqueue(callerAddress []byte, startNode *TestCallNode, callType vm.CallType,
	asyncData []byte, data []byte) {
	callToEnqueue := createCrossShardCall(startNode, callerAddress, callType, asyncData, data)
	queue.Data = append(queue.Data, callToEnqueue)
	sort.Stable(queue)
}

func createCrossShardCall(startNode *TestCallNode, callerAddress []byte, callType vm.CallType,
	asyncData []byte, data []byte) *CrossShardCall {
	parentsPath := make([]*TestCallNode, 0)
	crtNode := startNode
	for crtNode.Parent != nil {
		// we add parents for the enqueued node until we reach a callback edge
		if crtNode.Parent.IsCallback() {
			break
		}
		parentsPath = append(parentsPath, crtNode.Parent)
		crtNode = crtNode.Parent
	}
	callToEnqueue := &CrossShardCall{
		CallerAddress: callerAddress,
		StartNode:     startNode,
		CallType:      callType,
		AsyncData:     asyncData,
		Data:          data,
		ParentsPath:   parentsPath,
	}
	return callToEnqueue
}

// Requeue -
func (queue *CrossShardCallsQueue) Requeue(crossShardCall *CrossShardCall) {
	queue.Enqueue(crossShardCall.CallerAddress, crossShardCall.StartNode, crossShardCall.CallType,
		crossShardCall.AsyncData, crossShardCall.Data)
}

// Top -
func (queue *CrossShardCallsQueue) Top() *CrossShardCall {
	if len(queue.Data) == 0 {
		return nil
	}
	return queue.Data[0]
}

// Dequeue -
func (queue *CrossShardCallsQueue) Dequeue() *CrossShardCall {
	top := queue.Top()
	if top == nil {
		return nil
	}
	queue.Data = queue.Data[1:]
	return top
}

// IsEmpty -
func (queue *CrossShardCallsQueue) IsEmpty() bool {
	return len(queue.Data) == 0
}

// CanExecuteLocalCallback - in case of async local calls, search queue for pending children of the start of this edge
func (queue *CrossShardCallsQueue) CanExecuteLocalCallback(callbackNode *TestCallNode) bool {
	for _, callInQueue := range queue.Data {
		for _, parentInPath := range callInQueue.ParentsPath {
			if parentInPath == callbackNode.Parent {
				return false
			}
		}
	}
	return true
}

// Len -
func (queue *CrossShardCallsQueue) Len() int { return len(queue.Data) }

// Less -
func (queue *CrossShardCallsQueue) Less(i, j int) bool {
	return queue.Data[i].StartNode.ExecutionRound < queue.Data[j].StartNode.ExecutionRound
}

// Swap -
func (queue *CrossShardCallsQueue) Swap(i, j int) {
	queue.Data[i], queue.Data[j] = queue.Data[j], queue.Data[i]
}
