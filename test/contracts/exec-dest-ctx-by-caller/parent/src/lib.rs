#![no_std]
#![no_main]
#![allow(unused_attributes)]
#![feature(lang_items)]

use elrond_wasm::ContractIOApi;
use elrond_wasm_node::ArwenApiImpl;

const ALLOCSIZE: usize = 256;
const CHILDADDRESS: &[u8; 32] = b"\0\0\0\0\0\0\0\0\x0F\x0FchildSC...............";

pub static EEI: ArwenApiImpl = ArwenApiImpl{};

#[no_mangle]
pub extern "C" fn call_child() {
    let value: [u8; 32] = [0; 32];
    let arguments: &[&[u8]] = &[&[42 as u8]];
    execute_on_dest_context_by_caller(CHILDADDRESS, &value, "give", &arguments, 800000);
    EEI.finish_slice_u8(b"child called");
}

extern {
    fn executeOnDestContextByCaller(
        gas: i64,
        addressOffset: *const u8,
        valueOffset: *const u8,
        functionOffset: *const u8,
        functionLength: i32,
        numArguments: i32,
        argumentsLengthsOffset: *const i32,
        dataOffset: *const u8,
    ) -> i32;
}

pub fn execute_on_dest_context_by_caller(
    address: &[u8],
    value: &[u8],
    function: &str,
    arguments: &[&[u8]],
    gas: i64,
) -> i32
{
    let arguments_buffer = join_small_slices(arguments);
    let arguments_lengths = get_slices_lengths(arguments);

    unsafe {
        executeOnDestContextByCaller(
            gas,
            address.as_ref().as_ptr(),
            value.as_ref().as_ptr(),
            function.as_ptr(),
            function.len() as i32,
            arguments.len() as i32,
            arguments_lengths.as_ref().as_ptr(),
            arguments_buffer.as_ref().as_ptr(),
        )
    }
}

pub fn join_small_slices(slices: &[&[u8]]) -> [u8; ALLOCSIZE] {
    let mut buffer: [u8; ALLOCSIZE] = [0; ALLOCSIZE];
    let mut offset: usize = 0;

    for slice in slices {
        for byte in *slice {
            buffer[offset] = *byte;
            offset += 1;
        }
    }
    
    buffer
}

pub fn get_slices_lengths(slices: &[&[u8]]) -> [i32; ALLOCSIZE] {
    let mut lengths: [i32; ALLOCSIZE] = [0; ALLOCSIZE];
    for i in 0..slices.len() {
        lengths[i] = slices[i].len() as i32;
    }

    lengths
}
