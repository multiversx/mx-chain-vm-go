# WARNING: THIS CODE GENERATOR IS DEPRECATED
# TODO: Reimplement using go code generators if possible

from argparse import ArgumentParser


class HookSignature:
    def __init__(self, name, input, output, error=False, badReturn="nil"):
        super().__init__()
        self.name = name
        self.input = input
        self.output = output
        self.error = error
        self.badReturn = badReturn


signatures = [
    HookSignature(name="AccountExists", input=[("address", "[]byte")], output=[("result", "bool")], error=True, badReturn="false"),
    HookSignature(name="NewAddress", input=[("creatorAddress", "[]byte"), ("creatorNonce", "uint64"), ("vmType", "[]byte")], output=[("result", "[]byte")], error=True),
    HookSignature(name="GetBalance", input=[("address", "[]byte")], output=[("balance", "*big.Int")], error=True),
    HookSignature(name="GetNonce", input=[("address", "[]byte")], output=[("nonce", "uint64")], error=True, badReturn="0"),
    HookSignature(name="GetStorageData", input=[("address", "[]byte"), ("index", "[]byte")], output=[("data", "[]byte")], error=True),
    HookSignature(name="IsCodeEmpty", input=[("address", "[]byte")], output=[("result", "bool")], error=True, badReturn="false"),
    HookSignature(name="GetCode", input=[("address", "[]byte")], output=[("code", "[]byte")], error=True),
    HookSignature(name="GetBlockhash", input=[("nonce", "uint64")], output=[("result", "[]byte")], error=True),

    HookSignature(name="LastNonce", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="LastRound", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="LastTimeStamp", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="LastRandomSeed", input=[], output=[("result", "[]byte")], badReturn="nil"),
    HookSignature(name="LastEpoch", input=[], output=[("result", "uint32")], badReturn="0"),
    HookSignature(name="GetStateRootHash", input=[], output=[("result", "[]byte")], badReturn="nil"),
    HookSignature(name="CurrentNonce", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="CurrentRound", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="CurrentTimeStamp", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="CurrentRandomSeed", input=[], output=[("result", "[]byte")], badReturn="nil"),
    HookSignature(name="CurrentEpoch", input=[], output=[("result", "uint32")], badReturn="0")
]


def main():
    parser = ArgumentParser()
    subparsers = parser.add_subparsers()

    messages_parser = subparsers.add_parser("messages")
    messages_parser.set_defaults(func=generate_messages)

    repliers_parser = subparsers.add_parser("repliers")
    repliers_parser.set_defaults(func=generate_repliers)

    reply_slots_parser = subparsers.add_parser("reply-slots")
    reply_slots_parser.set_defaults(func=generate_reply_slots)

    gateway_parser = subparsers.add_parser("gateway")
    gateway_parser.set_defaults(func=generate_gateway)

    factory_parser = subparsers.add_parser("factory")
    factory_parser.set_defaults(func=generate_factory)

    args = parser.parse_args()

    if not hasattr(args, "func"):
        parser.print_help()
    else:
        args.func(args)


def generate_messages(args):
    print("package common")

    for signature in signatures:
        request_kind = f"Blockchain{signature.name}Request"
        request_go = f"""
        // Message{request_kind} represents a request message
        type Message{request_kind} struct {{
            Message
            {get_struct_fields_go(signature.input)}
        }}

        // NewMessage{request_kind} creates a request message
        func NewMessage{request_kind}({get_ctor_args(signature.input)}) *Message{request_kind} {{
            message := &Message{request_kind}{{}}
            message.Kind = {request_kind}
            {get_field_assignments(signature.input)}
            return message
        }}
        """

        response_kind = f"Blockchain{signature.name}Response"
        response_go = f"""
        // Message{response_kind} represents a response message
        type Message{response_kind} struct {{
            Message
             {get_struct_fields_go(signature.output)}
        }}

        // NewMessage{response_kind} creates a response message
        func NewMessage{response_kind}({get_ctor_args(signature.output, error=signature.error)}) *Message{response_kind} {{
            message := &Message{response_kind}{{}}
            message.Kind = {response_kind}
            {get_field_assignments(signature.output, error=signature.error)}
            return message
        }}
        """

        print(request_go)
        print(response_go)


def get_struct_fields_go(input_output):
    fields = []
    for arg_name, arg_type in input_output:
        field_name = my_capitalize(arg_name)
        fields.append(f"{field_name} {arg_type}")

    return "\n".join(fields)


def get_ctor_args(input_output, error=False):
    args = []
    for arg_name, arg_type in input_output:
        args.append(f"{arg_name} {arg_type}")

    if error:
        args.append(f"err error")

    return ", ".join(args)


def get_field_assignments(input_output, error=False):
    assignments = []
    for arg_name, arg_type in input_output:
        field_name = my_capitalize(arg_name)
        assignments.append(f"message.{field_name} = {arg_name}")

    if error:
        assignments.append(f"message.SetError(err)")

    return "\n".join(assignments)


def my_capitalize(input):
    return input[0].upper() + input[1:]


def generate_repliers(args):
    print("package nodepart")
    print("import \"github.com/ElrondNetwork/arwen-wasm-vm/v1.3/ipc/common\"")

    for signature in signatures:
        call_go, output_args = get_call(signature)
        typedRequest = f"typedRequest := request.(*common.MessageBlockchain{signature.name}Request)\n" if signature.input else ""

        func_go = f"""
        func (part *NodePart) replyToBlockchain{signature.name}(request common.MessageHandler) common.MessageHandler {{
            {typedRequest}{call_go}
            response := common.NewMessageBlockchain{signature.name}Response({output_args})
            return response
        }}
        """
        print(func_go)


def get_call(signature):
    output_args = []
    call_args = []

    for arg_name, _ in signature.output:
        output_args.append(arg_name)

    if signature.error:
        output_args.append(f"err")

    for arg_name, _ in signature.input:
        call_args.append(f"typedRequest.{my_capitalize(arg_name)}")

    output_args = ", ".join(output_args)
    call_args = ", ".join(call_args)

    return f"{output_args} := part.blockchain.{signature.name}({call_args})", output_args


def generate_reply_slots(args):
    print("part.Repliers = common.CreateReplySlots()")

    for signature in signatures:
        print(f"part.Repliers[common.Blockchain{signature.name}Request] = part.replyToBlockchain{signature.name}")


def generate_gateway(args):
    print("""
package arwenpart

import (
    "math/big"

    "github.com/ElrondNetwork/arwen-wasm-vm/v1.3/ipc/common"
    "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookGateway)(nil)

// BlockchainHookGateway forwards requests to the actual hook
type BlockchainHookGateway struct {
    messenger *ArwenMessenger
}

// NewBlockchainHookGateway creates a new gateway
func NewBlockchainHookGateway(messenger *ArwenMessenger) *BlockchainHookGateway {
    return &BlockchainHookGateway{messenger: messenger}
}
""")

    for signature in signatures:
        func_go = f"""
        // {signature.name} forwards a message to the actual hook
        func (blockchain *BlockchainHookGateway) {signature.name}({get_ctor_args(signature.input)}) {get_output_args(signature)} {{
            request := common.NewMessageBlockchain{signature.name}Request({get_call_args(signature)})
            rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
            if err != nil {{
                return {signature.badReturn}{", err" if signature.error else ""}
            }}

            if rawResponse.GetKind() != common.Blockchain{signature.name}Response {{
                return {signature.badReturn}{", common.ErrBadHookResponseFromNode" if signature.error else ""}
            }}

            response := rawResponse.(*common.MessageBlockchain{signature.name}Response)
            {get_gateway_return(signature)}
        }}
        """

        print(func_go)


def get_call_args(signature):
    call_args = []
    for arg_name, _ in signature.input:
        call_args.append(arg_name)

    return ", ".join(call_args)


def get_output_args(signature):
    output_args = []
    for arg_name, arg_type in signature.output:
        output_args.append(arg_type)
    if signature.error:
        output_args.append(f"error")

    result = ", ".join(output_args)
    if len(output_args) > 1:
        result = f"({result})"

    return result


def get_gateway_return(signature):
    result_field, _ = signature.output[0]
    result_field = my_capitalize(result_field)

    if signature.error:
        return f"return response.{result_field}, response.GetError()"
    return f"return response.{result_field}"


def generate_factory(args):
    assignments = ""

    for signature in signatures:
        assignments += f"messageCreators[Blockchain{signature.name}Request] = createMessageBlockchain{signature.name}Request"
        assignments += "\n"
        assignments += f"messageCreators[Blockchain{signature.name}Response] = createMessageBlockchain{signature.name}Response"
        assignments += "\n"

    print(f"""
package common


// CreateMessage creates a message given its kind
func CreateMessage(kind MessageKind) MessageHandler {{
    kindIndex := uint32(kind)
    length := uint32(len(messageCreators))
    if kindIndex < length {{
        message := messageCreators[kindIndex]()
        message.SetKind(kind)
        return message
    }}

    return createUndefinedMessage()
}}

type messageCreator func() MessageHandler

var messageCreators = make([]messageCreator, LastKind)

func init() {{
    for i := 0; i < len(messageCreators); i++ {{
        messageCreators[i] = createUndefinedMessage
    }}

    messageCreators[Initialize] = createMessageInitialize
    messageCreators[Stop] = createMessageStop
    messageCreators[ContractDeployRequest] = createMessageContractDeployRequest
    messageCreators[ContractCallRequest] = createMessageContractCallRequest
    messageCreators[ContractResponse] = createMessageContractResponse
    messageCreators[DiagnoseWaitRequest] = createMessageDiagnoseWaitRequest
    messageCreators[DiagnoseWaitResponse] = createMessageDiagnoseWaitResponse

    {assignments}
}}

func createMessageInitialize() MessageHandler {{
    return &MessageInitialize{{}}
}}

func createMessageStop() MessageHandler {{
    return &MessageStop{{}}
}}

func createMessageContractDeployRequest() MessageHandler {{
    return &MessageContractDeployRequest{{}}
}}

func createMessageContractCallRequest() MessageHandler {{
    return &MessageContractCallRequest{{}}
}}

func createMessageContractResponse() MessageHandler {{
    return &MessageContractResponse{{}}
}}

func createMessageDiagnoseWaitRequest() MessageHandler {{
    return &MessageDiagnoseWaitRequest{{}}
}}

func createMessageDiagnoseWaitResponse() MessageHandler {{
    return &MessageDiagnoseWaitResponse{{}}
}}

func createUndefinedMessage() MessageHandler {{
    return NewUndefinedMessage()
}}
""")

    for signature in signatures:
        print(f"""
        func createMessageBlockchain{signature.name}Request() MessageHandler {{
            return &MessageBlockchain{signature.name}Request{{}}
        }}
        """)

        print(f"""
        func createMessageBlockchain{signature.name}Response() MessageHandler {{
            return &MessageBlockchain{signature.name}Response{{}}
        }}
        """)


if __name__ == "__main__":
    main()
