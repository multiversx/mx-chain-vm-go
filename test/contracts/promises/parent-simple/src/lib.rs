#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::api::{EndpointArgumentApi, EndpointFinishApi, StorageReadApi, StorageWriteApi};
use elrond_wasm_node::ArwenApiImpl;

use promises_common::*;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

const SUCCESS_CALLBACK_ARGUMENT_KEY: &[u8] = b"SuccessCallbackArg";
const FAIL_CALLBACK_ARGUMENT_KEY: &[u8] = b"FailCallbackArg";
const CURRENT_STORAGE_INDEX_KEY: &[u8] = b"CurrentStorageIndex";

const COMMON_GROUP_ID: &[u8] = b"testgroup";
const SUCCESS_CALLBACK_ONE_ARG_NAME: &[u8] = b"success_callback_one_arg";
const FAIL_CALLBACK_NAME: &[u8] = b"fail_callback";

#[no_mangle]
pub extern "C" fn no_async() {
    EEI.finish_i64(42);
}

// one async call

#[no_mangle]
pub extern "C" fn one_async_call_no_cb_with_call_value() {
    let mut value = [0u8; 32];
    value[31] = 16;

    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &value,
                      b"answer",
                      EMPTY_SLICE,
                      EMPTY_SLICE,
                      100000);
}

#[no_mangle]
pub extern "C" fn one_async_call_no_cb_fail() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"fail",
                      EMPTY_SLICE,
                      EMPTY_SLICE,
                      100000);
}

#[no_mangle]
pub extern "C" fn one_async_call_no_cb_fail_with_call_value() {
    let mut value = [0u8; 32];
    value[31] = 16;

    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &value,
                      b"fail",
                      EMPTY_SLICE,
                      EMPTY_SLICE,
                      100000);
}

#[no_mangle]
pub extern "C" fn one_async_call_success_cb() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"answer",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      EMPTY_SLICE,
                      100000);
}

#[no_mangle]
pub extern "C" fn one_async_call_fail_cb() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"fail",
                      EMPTY_SLICE,
                      FAIL_CALLBACK_NAME,
                      100000);
}

// two async calls

pub extern "C" fn two_async_same_cb_success_both() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"echo@0x01",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
    
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"echo@0x02",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
}

pub extern "C" fn two_async_same_cb_success_first_fail_second() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"echo@0x01",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
    
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"fail",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
}

pub extern "C" fn two_async_same_cb_fail_first_success_second() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"fail",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
    
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"echo@0x02",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
}

pub extern "C" fn two_async_same_cb_fail_both() {
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"fail",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
    
    create_async_call(COMMON_GROUP_ID,
                      &Address::from(CHILD_ADDRESS),
                      &ZERO,
                      b"fail",
                      SUCCESS_CALLBACK_ONE_ARG_NAME,
                      FAIL_CALLBACK_NAME,
                      100000);
}

// callbacks

// first argument is "0" for success, followed by data passed by finish() in callee contract
#[no_mangle]
pub extern "C" fn success_callback_one_arg() {
    let expected_num_args = 2;
    EEI.check_num_arguments(expected_num_args);

    let mut storage_index = EEI.storage_load_u64(&CURRENT_STORAGE_INDEX_KEY);

    for arg_index in 0..expected_num_args {
        let arg = EEI.get_argument_u64(arg_index);
        let storage_key = construct_storage_key(&[SUCCESS_CALLBACK_ARGUMENT_KEY, &[storage_index as u8]]);

        storage_index += 1;
        EEI.storage_store_u64(&storage_key, arg);
    }

    EEI.storage_store_u64(&CURRENT_STORAGE_INDEX_KEY, storage_index);
}

// first argument is error code, followed by error message
#[no_mangle]
pub extern "C" fn fail_callback() {
    let expected_num_args = 2;
    EEI.check_num_arguments(expected_num_args);

    let mut storage_index = EEI.storage_load_u64(&CURRENT_STORAGE_INDEX_KEY);

    for arg_index in 0..expected_num_args {
        let arg = EEI.get_argument_vec_u8(arg_index);
        let storage_key = construct_storage_key(&[FAIL_CALLBACK_ARGUMENT_KEY, &[storage_index as u8]]);
    
        storage_index += 1;
        EEI.storage_store_slice_u8(&storage_key, &arg);
    }

    EEI.storage_store_u64(&CURRENT_STORAGE_INDEX_KEY, storage_index);
}

