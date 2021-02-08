#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::{ContractIOApi, Address};
use elrond_wasm_node::ArwenApiImpl;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

#[no_mangle]
pub extern "C" fn answer() {
    EEI.finish_i64(42);
}
