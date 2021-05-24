ALICE="/home/elrond/Downloads/TO_MOVE/TO_MOVE/WalletKey.pem"
ADDRESS=$(erdpy data load --key=address-testnet)
DEPLOY_TRANSACTION=$(erdpy data load --key=deployTransaction-testnet)
PROXY=https://testnet-api.elrond.com
PROJECT="/home/elrond/Github/revert_esdt_test/revert-esdt-test/"

deploy() {
    erdpy --verbose contract deploy --project=${PROJECT} --recall-nonce --pem=${ALICE} --gas-limit=50000000 --send --outfile="deploy-testnet.interaction.json" --metadata-payable --proxy=${PROXY} --chain=T || return

    TRANSACTION=$(erdpy data parse --file="deploy-testnet.interaction.json" --expression="data['emitted_tx']['hash']")
    ADDRESS=$(erdpy data parse --file="deploy-testnet.interaction.json" --expression="data['emitted_tx']['address']")

    erdpy data store --key=address-testnet --value=${ADDRESS}
    erdpy data store --key=deployTransaction-testnet --value=${TRANSACTION}

    echo ""
    echo "Smart contract address: ${ADDRESS}"
}

upgrade() {
    erdpy --verbose contract upgrade erd1qqqqqqqqqqqqqpgqczgetjuzakt3ug9sr2gfd7kcdcspuk03t9usm4m54k --project=${PROJECT} --recall-nonce --pem=${ALICE} --gas-limit=50000000 --send --outfile="deploy-testnet.interaction.json" --metadata-payable --proxy=${PROXY} --chain=T
    sleep 6
    erdpy --verbose contract upgrade erd1qqqqqqqqqqqqqpgqaupztl0yyv77kqgj52pk9h9rmm8seh87t9usny723u --project=${PROJECT} --recall-nonce --pem=${ALICE} --gas-limit=50000000 --send --outfile="deploy-testnet.interaction.json" --metadata-payable --proxy=${PROXY} --chain=T
    sleep 6
    erdpy --verbose contract upgrade erd1qqqqqqqqqqqqqpgq9wcy0mcu3u79pga3v3jttm9rpqkrrj84t9usq224ax --project=${PROJECT} --recall-nonce --pem=${ALICE} --gas-limit=50000000 --send --outfile="deploy-testnet.interaction.json" --metadata-payable --proxy=${PROXY} --chain=T
}

setBurnLocalRole() {
    erdpy --verbose contract call erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u --recall-nonce --pem=${ALICE} --gas-limit=1000000000 --send --function="setSpecialRole" --arguments 0x414c432d313764386632 0x00000000000000000500ef0225fde4233deb0112a28362dca3decf0cdcfe5979 0x45534454526f6c654c6f63616c4275726e --proxy=${PROXY} --chain=T
}

test() {
    erdpy --verbose contract call erd1qqqqqqqqqqqqqpgqczgetjuzakt3ug9sr2gfd7kcdcspuk03t9usm4m54k --recall-nonce --pem=${ALICE} --gas-limit=1000000000 --send --function="transferExecuteFungibleAndFailAfterBurn" --arguments 0x00000000000000000500ef0225fde4233deb0112a28362dca3decf0cdcfe5979 0x414c432d313764386632 0x100000 --proxy=${PROXY} --chain=T
}

test2() {
    erdpy --verbose contract call erd1qqqqqqqqqqqqqpgqczgetjuzakt3ug9sr2gfd7kcdcspuk03t9usm4m54k --recall-nonce --pem=${ALICE} --gas-limit=1000000000 --send --function="executeOnDestWithFungiblePaymentAndFailAfterBurn" --arguments 0x000000000000000005002bb047ef1c8f3c50a3b16464b5eca3082c31c8f55979 0x414c432d313764386632 0x100000 --proxy=${PROXY} --chain=T
}