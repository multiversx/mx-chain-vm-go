#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::{ContractHookApi, ContractIOApi, Address};
use elrond_wasm_node::ArwenApiImpl;
use elrond_wasm_node::ArwenBigUint;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

#[no_mangle]
pub extern "C" fn give() {
    let value_to_give = EEI.get_argument_u8(0);
    let caller: Address = EEI.get_caller();
    let value = ArwenBigUint::from(value_to_give as u64);
    EEI.send_tx(&caller, &value, "");
    EEI.finish_slice_u8(b"sent");
}
