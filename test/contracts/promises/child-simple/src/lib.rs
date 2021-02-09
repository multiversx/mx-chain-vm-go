#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::api::{EndpointFinishApi, ErrorApi};
use elrond_wasm_node::ArwenApiImpl;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

#[no_mangle]
pub extern "C" fn answer() {
    EEI.finish_u64(42);
}

#[no_mangle]
pub extern "C" fn fail() {
    EEI.signal_error(&b"fail"[..]);
}
