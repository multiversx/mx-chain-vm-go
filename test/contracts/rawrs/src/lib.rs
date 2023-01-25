#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]




#[no_mangle]
pub extern "C" fn method() {
    let api = ArwenApiImpl{};
    api.finish_i64(4);
}
