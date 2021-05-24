#![no_std]
#![allow(non_snake_case)]

elrond_wasm::imports!();
elrond_wasm::derive_imports!();

mod ProxyMod {
    elrond_wasm::imports!();

    #[elrond_wasm_derive::proxy]
    pub trait ProxySelf {
        #[payable("*")]
        #[endpoint]
        fn burnAndFail(&self);
    }
}

#[elrond_wasm_derive::contract]
pub trait RevertEsdtTester {
    #[proxy]
    fn my_proxy(&self, to: Address) -> ProxyMod::Proxy<Self::SendApi>;

    #[init]
    fn init(&self) {}

    // This is not ok, Destination will lose balance while revertEsdt will try to send back the funds.
    #[endpoint]
    fn transferExecuteFungibleAndFailAfterBurn(
        &self,
        address: Address,
        token_id: TokenIdentifier,
        amount: Self::BigUint,
    ) -> SCResult<(Self::BigUint, Self::BigUint, bool)> {
        let balance_before = self.blockchain().get_esdt_balance(
            &self.blockchain().get_sc_address(),
            token_id.as_esdt_identifier(),
            0,
        );

        let transf_exec_result = self.send().direct_esdt_execute(
            &address,
            token_id.as_esdt_identifier(),
            &amount,
            self.blockchain().get_gas_left() / 2,
            b"burnAndFail",
            &ArgBuffer::new(),
        );

        let transfer_success = match transf_exec_result {
            Result::Ok(()) => true,
            Result::Err(_) => false,
        };

        let balance_after = self.blockchain().get_esdt_balance(
            &self.blockchain().get_sc_address(),
            token_id.as_esdt_identifier(),
            0,
        );

        Ok((balance_before, balance_after, transfer_success))
    }

    // This is ok, entire transaction fails. No moves of balance.
    #[endpoint]
    fn executeOnDestWithFungiblePaymentAndFailAfterBurn(
        &self,
        address: Address,
        token_id: TokenIdentifier,
        amount: Self::BigUint,
    ) -> SCResult<(Self::BigUint, Self::BigUint)> {
        let balance_before = self.blockchain().get_esdt_balance(
            &self.blockchain().get_sc_address(),
            token_id.as_esdt_identifier(),
            0,
        );

        self.my_proxy(address)
            .burnAndFail()
            .with_token_transfer(token_id.clone(), amount.clone())
            .execute_on_dest_context(self.blockchain().get_gas_left());

        let balance_after = self.blockchain().get_esdt_balance(
            &self.blockchain().get_sc_address(),
            token_id.as_esdt_identifier(),
            0,
        );

        Ok((balance_before, balance_after))
    }

    #[payable("*")]
    #[endpoint]
    fn burnAndFail(&self) -> SCResult<()> {
        let (amount, token_id) = self.call_value().payment_token_pair();

        self.send()
            .burn_tokens(&token_id, 0, &amount, self.blockchain().get_gas_left() / 2);

        sc_error!("Burned tokens and returned Error")
    }
}
