#![no_std]

imports!();

const TOKEN_NAME: &[u8] = b"TT";


pub trait Exchange {

    #[endpoint(validateGetters)]
    fn validate_getters(&self) -> SCResult<()> {
        sc_try!(self.validate_esdt_token_name());
        sc_try!(self.validate_esdt_token_value(5));
        Ok(())
    }

    fn validate_esdt_token_name(&self) -> SCResult<()> {
        let token_name: Option<Vec<u8>> = self.get_esdt_token_name();
        match token_name {
            None => {
                sc_error!("esdt token required")
            },
            Some(name) => {
                require!(name.as_slice() == TOKEN_NAME, "wrong esdt token");
                Ok(())
            }
        }
    }

    fn validate_esdt_token_value(&self, expected_value: u64) -> SCResult<()> {
        let token_value = self.get_esdt_value_big_uint();
        let expected_value = BigUint::from(expected_value);
        require!(expected_value == token_value, "wrong esdt value");
        Ok(())
    }

    #[endpoint(validateGettersAfterESDTTransfer)]
    fn validate_getters_after_esdt_transfer(&self) -> SCResult<()> {
        Ok(())
    }
}
