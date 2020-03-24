package common

// MessageDiagnoseWaitRequest is a diagnose request message (from Node)
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

// MessageDiagnoseWaitResponse is a diagnose response message (from Arwen)
type MessageDiagnoseWaitResponse struct {
	Message
}

// NewMessageDiagnoseWaitResponse creates a message
func NewMessageDiagnoseWaitResponse() *MessageDiagnoseWaitResponse {
	message := &MessageDiagnoseWaitResponse{}
	message.Kind = DiagnoseWaitResponse
	return message
}
