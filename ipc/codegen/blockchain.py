
import collections
from argparse import ArgumentParser

HookSignature = collections.namedtuple("HookSignature", ["name", "input", "output", "error"], verbose=False, rename=False)

signatures = [
    HookSignature(name="AccountExists", input=[("address", "[]byte")], output=[("result", "bool")], error=True),
    HookSignature(name="NewAddress", input=[("creatorAddress", "[]byte"), ("creatorNonce", "uint64"), ("vmType", "[]byte")], output=[("result", "[]byte")], error=True),
    HookSignature(name="GetBalance", input=[("address", "[]byte")], output=[("balance", "*big.Int")], error=True),
    HookSignature(name="GetNonce", input=[("address", "[]byte")], output=[("nonce", "uint64")], error=True),
    HookSignature(name="GetStorageData", input=[("address", "[]byte"), ("index", "[]byte")], output=[("data", "[]byte")], error=True),
    HookSignature(name="IsCodeEmpty", input=[("address", "[]byte")], output=[("result", "bool")], error=True),
    HookSignature(name="GetCode", input=[("address", "[]byte")], output=[("code", "[]byte")], error=True),
    HookSignature(name="GetBlockhash", input=[("nonce", "uint64")], output=[("result", "[]byte")], error=True),

    HookSignature(name="LastNonce", input=[], output=[("result", "uint64")], error=False),
    HookSignature(name="LastRound", input=[], output=[("result", "uint64")], error=False),
    HookSignature(name="LastTimeStamp", input=[], output=[("result", "uint64")], error=False),
    HookSignature(name="LastRandomSeed", input=[], output=[("result", "[]byte")], error=False),
    HookSignature(name="LastEpoch", input=[], output=[("result", "uint32")], error=False),
    HookSignature(name="GetStateRootHash", input=[], output=[("result", "[]byte")], error=False),
    HookSignature(name="CurrentNonce", input=[], output=[("result", "uint64")], error=False),
    HookSignature(name="CurrentRound", input=[], output=[("result", "uint64")], error=False),
    HookSignature(name="CurrentTimeStamp", input=[], output=[("result", "uint64")], error=False),
    HookSignature(name="CurrentRandomSeed", input=[], output=[("result", "[]byte")], error=False),
    HookSignature(name="CurrentEpoch", input=[], output=[("result", "uint32")], error=False)
]


def main():
    parser = ArgumentParser()
    subparsers = parser.add_subparsers()

    messages_parser = subparsers.add_parser("messages")
    messages_parser.set_defaults(func=generate_messages)
    replies_parser = subparsers.add_parser("replies")
    replies_parser.set_defaults(func=generate_replies)

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


def generate_replies(args):
    print("package nodepart")
    print("import \"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common\"")

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


if __name__ == "__main__":
    main()