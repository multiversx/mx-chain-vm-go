# arwen-wasm-vm

[![Build Status](https://travis-ci.com/ElrondNetwork/arwen-wasm-vm.svg?branch=master)](https://travis-ci.com/ElrondNetwork/arwen-wasm-vm)

Arwen is the WASM-based Virtual Machine for running Elrond Smart Contracts.

The VM is launched as a child process of the [Elrond Node](https://github.com/ElrondNetwork/elrond-go).

## IPC communication with the Node

The parent process communicates with Arwen by means of pipes (in-memory `os.File` objects)
. The main components involved in the communication flow are: 
 - **Arwen Driver** - used by Node to manage Arwen's process)
 - **Node Part** - contains Node's main loop of messages; the node starts a message loop for each contract request (deployment / execution)
 - **Arwen Part** - Arwen's forever loop of messages


## API Hooks

### CryptoHook 

This hook is implemented directly in Arwen. There is no cross-process communication involved in executing functions of the `CryptoHook`.

### BlockchainHook

Any function call against the `BlockchainHook`, which may happen during the execution of a smart contract results in the following:

 1. The function call is forwarded to `BlockchainGateway`
 1. `BlockchainGateway` sends a request message to the Node, via a `Messenger`, through the pipe `arwenToNode`.
 1. `BlockchainGateway` waits indefinitely until a response message comes from the Node through the pipe `nodeToArwen`.
 1. `NodePart` receives the request message from Arwen, via its own `Messenger` and resolves the request against the actual `BlockchainHook` implementation (within the appropriate **blockchain replier**).
 1. The result of the `BlockchainHook` call is wrapped in a message and sent to Arwen, through the pipe `nodeToArwen`.
 1. The control returns to `BlockchainGateway` and then to the smart contract.


## Messaging

#### Pipes - the communication takes place through these pipes:

 - `arwenInit` - Arwen initialization parameters (arguments) are passed through this pipe.
 - `nodeToArwen` - used to transport contract requests (deployment and execution) towards Arwen and also to respond on Arwen’s blockchain hook calls.
 - `arwenToNode` - used to transport blockchain hook calls from Arwen to the Node and also to respond with contract deployment / execution results.

#### Dialogue

A dialogue between the Node and Arwen starts with the deployment or execution of a smart contract and consists of a sequence of messages. 

Before a `Messenger` component sends a `Message`, it labels it with a **dialogue nonce** (an increasing integer). The dialogue ends (resets) when the deployment or execution is finished (the nonce is also reset at this very time).

Both Arwen’s Messenger and Node’s Messenger check the dialogue nonces of the incoming messages to ensure the correctness of the dialogue (safety net).

#### Data protocol

When a Messenger needs to send a message, it first sends a preamble, and then the payload.

When a Messenger needs to receive a message, it first reads the preamble from the pipe, and then the payload.  

The preamble consists in:

 - The length of the payload (4 bytes integer)
 - The kind of the message to be sent (4 bytes integer, corresponding to the type `MessageKind`).

#### Serialization

Currently, JSON format is used to format the messages before sending them through the pipe.

#### Loops

The Arwen Part contains an **infinite message loop**. When the loop is broken (in case of a critical error), Arwen stops. The Node Part starts a message loop for each contract request it needs to send to Arwen. The message loop ends when the response is received (or in case of a critical error).


#### Blocking reads

`Messenger` components perform blocking reads against the pipes.  Read timeout is set by means of `SetDeadline` calls. See `Receiver`.

### Critical errors

Caught critical errors end Arwen’s main loop and Arwen’s process - the process will be restarted on the very next contract request - see `RestartArwenIfNecessary`.

Panics in Arwen’s process lead to a restart - performed by `Arwen Driver`.


### Path of Arwen’s binary

Arwen Driver will first look for the Arwen binary in Node’s current directory. If the binary isn’t found, it will look at the path specified by the environment variable `ARWEN_PATH`.


### Loggers

Logs are sent from Arwen to the Node through pipes as well. The Arwen Driver also captures Arwen’s `STDOUT` and `STDERR`.

Loggers defined on Arwen's part:

 - `arwen/host` 
 - `arwen/part`
 - `arwen/baseMessenger`
 - `arwen/duration`

Loggers defined on Node's part:

 - `arwenDriver`
 - `arwen/baseMessenger`
