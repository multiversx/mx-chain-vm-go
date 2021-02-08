#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::Address;
use elrond_wasm::api::EndpointFinishApi;
use elrond_wasm_node::ArwenApiImpl;

use promises_common::{
    create_async_call,
    CHILD_ADDRESS
};

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

#[no_mangle]
pub extern "C" fn no_async() {
    EEI.finish_i64(42);
}

#[no_mangle]
pub extern "C" fn one_async_call_no_cb() {
    let mut value: [u8; 32] = [0; 32];

    value[31] = 16;

    create_async_call("testgroup",
                      &Address::from(CHILD_ADDRESS),
                      &value,
                      b"answer",
                      "",
                      "",
                      100000);
}
