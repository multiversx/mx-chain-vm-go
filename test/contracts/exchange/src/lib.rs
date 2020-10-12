#![no_std]

imports!();

#[elrond_wasm_derive::contract(ExchangeImpl)]
pub trait Exchange {

    #[endpoint(acceptDonation)]
    fn accept_donation(&self) -> SCResult<()> {
        Ok(())
    }
}
