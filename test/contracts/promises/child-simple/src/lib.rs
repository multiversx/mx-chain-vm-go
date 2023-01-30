#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]




pub static EEI: VMHooksImpl = VMHooksImpl{};

#[no_mangle]
pub extern "C" fn answer() {
    EEI.finish_u64(42);
}

#[no_mangle]
pub extern "C" fn answer_wrong() {
    EEI.finish_u64(24);
}

// receives u64 as argument and returns it back
#[no_mangle]
pub extern "C" fn echo() {
    EEI.check_num_arguments(1);

    let arg = EEI.get_argument_u64(0);

    EEI.finish_u64(arg);
}

#[no_mangle]
pub extern "C" fn fail() {
    EEI.signal_error(&b"fail"[..]);
}
