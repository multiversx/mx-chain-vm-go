#![no_std]

pub use elrond_wasm::{Address, Vec};

pub const PARENT_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b'p', b'a', b'r', b'e', b'n', b't',
    b'S', b'C', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];

pub const CHILD_ADDRESS: [u8; 32] = [
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x0F, b'c', b'h', b'i', b'l', b'd', b'S',
    b'C', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.', b'.',
];

pub const ZERO: [u8; 32] = [0u8; 32];

pub fn construct_storage_key(key_parts: &[&[u8]]) -> Vec<u8> {
    let mut key = Vec::new();

    for part in key_parts {
        key.extend_from_slice(part);
    }

    key
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
