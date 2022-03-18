# Onchain Setup
These steps require familarity with cmd line and [hardhat](https://hardhat.org/getting-started/).

If BAS is available at localhost:8545, use `localhost` for `<network>` below. Otherwise, add proper network to hardhat config and use configured BAS network name.
## Preparation
- BEP20 contract on BAS as *cBridgeStakeToken*
- clone https://github.com/celer-network/sgn-v2-contracts/ repo
- set DEFAULT_PRIVATE_KEY to deployer account's private key in hardhat.config.ts

## Deploy SGN core contracts
- copy [.env.example](./.env.example) to `sgn-v2-contracts/.env`
- in .env file, set CELR to *cBridgeStakeToken* address on BAS. Update other parameters like MAX_VALIDATOR_NUM accordingly.
- run `hardhat deploy --network <network> --tags SGNStaking`
- record deployed contracts addresses for SGN config later

## Deploy cBridge contracts
First deploy on BAS:
- Deploy Bridge contract: `hardhat deploy --network <network>  --tags Bridge`
- in .env file, set `ORIGINAL_TOKEN_VAULT_SIGS_VERIFIER` and `PEGGED_TOKEN_BRIDGE_SIGS_VERIFIER` to the Bridge address from previous step
- Deploy OriginalTokenVault contract: `hardhat deploy --network <network>  --tags OriginalTokenVault`
- Deploy PeggedTokenBridge contract: `hardhat deploy --network <network>  --tags PeggedTokenBridge`

Repeat above steps for BSC, use `bsc` for `<network>`.

**CAUTION:** Make sure to update `ORIGINAL_TOKEN_VAULT_SIGS_VERIFIER` and `PEGGED_TOKEN_BRIDGE_SIGS_VERIFIER` in .env to BSC Bridge address before deploy OriginalTokenVault and PeggedTokenBridge.
