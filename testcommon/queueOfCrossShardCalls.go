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
}

// CrossShardCallsQueue -
type CrossShardCallsQueue struct {
	data []*CrossShardCall
}

// NewCrossShardCallQueue -
func NewCrossShardCallQueue() *CrossShardCallsQueue {
	return &CrossShardCallsQueue{}
}

// Enqueue -
func (queue *CrossShardCallsQueue) Enqueue(callerAddress []byte, startNode *TestCallNode, callType vm.CallType, data []byte) {
	queue.data = append(queue.data, &CrossShardCall{
		CallerAddress: callerAddress,
		StartNode:     startNode,
		CallType:      callType,
		Data:          data,
	})
}

// Top -
func (queue *CrossShardCallsQueue) Top() *CrossShardCall {
	if len(queue.data) == 0 {
		return nil
	}
	return queue.data[0]
}

// Dequeue -
func (queue *CrossShardCallsQueue) Dequeue() *CrossShardCall {
	top := queue.Top()
	if top == nil {
		return nil
	}
	queue.data = queue.data[1:]
	return top
}

// IsEmpty -
func (queue *CrossShardCallsQueue) IsEmpty() bool {
	return len(queue.data) == 0
}
