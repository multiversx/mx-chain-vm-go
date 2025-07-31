# Asynchronous Logic

This document explains the asynchronous logic used in the VM host.

## Overview

The asynchronous logic allows smart contracts to call other smart contracts without blocking the execution of the current smart contract. This is achieved by using a callback mechanism. When a smart contract calls another smart contract asynchronously, it provides a callback function that will be called when the asynchronous call is complete.

## Asynchronous Context

The `asyncContext` is responsible for handling asynchronous calls between smart contracts. It keeps track of the call stack, the arguments, and the return data.

The `asyncContext` has the following responsibilities:

- Initializing the internal state of the asynchronous context.
- Pushing and popping the state of the asynchronous context to and from the state stack.
- Getting and setting the return data.
- Getting and setting the call group.
- Registering and executing asynchronous calls.

## Asynchronous Call Groups

An asynchronous call group is a group of asynchronous calls that are executed together. This is useful when a smart contract needs to make multiple asynchronous calls and wants to be notified when all of them are complete.

## Asynchronous Calls

An asynchronous call is a call to a smart contract that is executed asynchronously. When an asynchronous call is made, the `asyncContext` creates a new `AsyncCall` object and adds it to the current `AsyncCallGroup`.

The `AsyncCall` object contains the following information:

- The destination address of the smart contract to be called.
- The function to be called.
- The arguments to be passed to the function.
- The value to be transferred to the smart contract.
- The gas limit for the call.
- The success and error callbacks.

When the asynchronous call is complete, the `asyncContext` calls the appropriate callback function.

## Legacy Asynchronous Calls

Legacy asynchronous calls are asynchronous calls that were made before the introduction of the `AsyncCall` object. Legacy asynchronous calls are still supported for backward compatibility.
