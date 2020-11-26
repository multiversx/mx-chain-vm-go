#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::{ContractIOApi, Address};
use elrond_wasm_node::ArwenApiImpl;

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

extern {
    fn createAsyncCall(
        groupIDOffset: *const u8, groupIDLength: i32,
        destOffset: *const u8,
        valueOffset: *const u8,
        dataOffset: *const u8,
        dataLength: i32,
        successCallbackNameOffset: *const u8, successCallbackNameLen: i32,
        errorCallbackNameOffset: *const u8, errorCallbackNameLen: i32,
        gas: i64
    );
}

pub fn create_async_call(
    group_id: &str,
    destination: &Address,
    value: &[u8],
    data: &[u8],
    success_callback_name: &str,
    error_callback_name: &str,
    gas: i64) 
{
    unsafe {
        createAsyncCall(
            group_id.as_ptr(),
            group_id.len() as i32,
            destination.as_ref().as_ptr(),
            value.as_ref().as_ptr(),
            data.as_ref().as_ptr(),
            data.len() as i32,
            success_callback_name.as_ptr(),
            success_callback_name.len() as i32,
            error_callback_name.as_ptr(),
            error_callback_name.len() as i32,
            gas
        )
    }
}


#[no_mangle]
pub extern "C" fn no_async() {
    EEI.finish_i64(42);
}

#[no_mangle]
pub extern "C" fn one_async_call_no_cb() {

}
