package common

// MessageDiagnoseWaitRequest represents a message
type MessageDiagnoseWaitRequest struct {
	Message
	Milliseconds uint32
}

// NewMessageDiagnoseWaitRequest creates a message
func NewMessageDiagnoseWaitRequest(milliseconds uint32) *MessageDiagnoseWaitRequest {
	message := &MessageDiagnoseWaitRequest{}
	message.Kind = DiagnoseWaitRequest
	message.Milliseconds = milliseconds
	return message
}

// MessageDiagnoseWaitResponse is
type MessageDiagnoseWaitResponse struct {
	Message
}

// NewMessageDiagnoseWaitResponse creates a message
func NewMessageDiagnoseWaitResponse() *MessageDiagnoseWaitResponse {
	message := &MessageDiagnoseWaitResponse{}
	message.Kind = DiagnoseWaitResponse
	return message
}
