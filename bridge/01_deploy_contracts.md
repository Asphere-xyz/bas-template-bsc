# Deploy SGN and cBridge contracts

These steps require familiarity with command line and [hardhat](https://hardhat.org/getting-started) for contract deployment.

## Preparations

1. Prepare a BEP20 contract with decimal 18 on BAS as `SGNStakeToken`. You can use the wrapped version of your native gas token or deploy a
new one.

2. Clone the [SGN V2 contracts](https://github.com/celer-network/sgn-v2-contracts) repo and run `yarn install`.

3. Copy [contract-configs/.env.example](contract-configs/.env.example) to `sgn-v2-contracts/.env`.

4. Following existing examples, add the BAS network to `networks` in `hardhat.config.ts`. Set `BAS_ENDPOINT` to your BAS RPC endpoint.
Set `BSC_PRIVATE_KEY` and `BAS_PRIVATE_KEY` to your deployer key.

## Deploy SGN staking contracts

1. In `.env`, set `CELR` to the `SGNStakeToken` address on BAS. Update other parameters like `MAX_VALIDATOR_NUM` according to your need.

2. Run

    ```sh
    hardhat deploy --network <network> --tags SGNStaking
    ```

    Record the deployed contracts addresses for configuring the SGN later.

## Deploy cBridge contracts

You need to deploy three cBridge contracts:

`Bridge` is the liquidity pool bridge contract that provides the functionality of a liquidity pool based bridge and serves as the on-chain signature verifier
for other types of bridges.

`OriginalTokenVault` is required when an original token on this chain needs to be bridged to another chain via the pegged token bridge.

`PeggedTokenBridge` is required when a pegged token for an original token on another chain needs to be bridged to the new chain via the pegged token bridge.

First, deploy on BAS:

1. Deploy `Bridge`:

    ```sh
    hardhat deploy --network <network>  --tags Bridge
    ```

    Record the deployed address.

2. In .env file, set `ORIGINAL_TOKEN_VAULT_SIGS_VERIFIER` and `PEGGED_TOKEN_BRIDGE_SIGS_VERIFIER` to the `Bridge` address from the previous
    step.

    Deploy `OriginalTokenVault`:

	```sh
	hardhat deploy --network <network>  --tags OriginalTokenVault
	```

    Deploy `PeggedTokenBridge`:

	```sh
	hardhat deploy --network <network>  --tags PeggedTokenBridge
	```

Repeat the steps above for BSC, using `bsc` for `<network>`. Record all deployed addresses.
