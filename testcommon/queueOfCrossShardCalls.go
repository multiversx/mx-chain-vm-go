package testcommon

import (
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
)

// CrossShardCall -
type CrossShardCall struct {
	CallerAddress []byte
	StartNode     *TestCallNode
	CallType      vm.CallType
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
func (queue *CrossShardCallsQueue) Enqueue(callerAddress []byte, startNode *TestCallNode, callType vm.CallType, data []byte) {
	parentsPath := make([]*TestCallNode, 0)
	crtNode := startNode
	for crtNode.Parent != nil {
		parentsPath = append(parentsPath, crtNode.Parent)
		crtNode = crtNode.Parent
	}
	queue.Data = append(queue.Data, &CrossShardCall{
		CallerAddress: callerAddress,
		StartNode:     startNode,
		CallType:      callType,
		Data:          data,
		ParentsPath:   parentsPath,
	})
}

// Requeue -
func (queue *CrossShardCallsQueue) Requeue(crossShardCall *CrossShardCall) {
	queue.Enqueue(crossShardCall.CallerAddress, crossShardCall.StartNode, crossShardCall.CallType, crossShardCall.Data)
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
