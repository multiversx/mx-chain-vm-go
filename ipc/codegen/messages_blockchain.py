import collections

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


if __name__ == "__main__":
    main()