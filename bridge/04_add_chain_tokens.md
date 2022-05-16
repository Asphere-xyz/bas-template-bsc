# Add chain and tokens

### Add new chain

These steps need to be repeated for your BAS and BSC respectively.

1. Edit the chain config in `cbridge.toml`. Example:

    ```toml
    [[multichain]]
    # The chain ID
    chainID = 560
    # The name of your BAS.
    name = "Your BAS"
    # The JSON-RPC endpoint URL.
    gateway = "https://<rpc-endpoint-url>"
    # The interval in seconds to poll on-chain events.
    blkinterval = 10
    # The number of block delays to consider an event finalized. Set as needed.
    blkdelay = 3
    # Some endpoints restrict the maximal number of blocks in eth_getLogs calls. Set as needed.
    maxblkdelta = 5000
    # Required to handle inconsistencies from the JSON-RPC providers. Set as needed.
    forwardblkdelay = 30
    # Required for timely confirmation of txs. Set as needed.
    addgasgwei = 1
    # Required in case gas estimation is off. Keep it as is.
    addgasestimateratio = 0.3
    cbridge = "<Bridge-contract-address>"
    otvault = "<OriginalTokenVault-contract-address>"
    ptbridge = "<PeggedTokenBridge-contract-address>"
    ```

2. Restart `sgnd` with the new config. Check for errors in `/var/log/sgnd/app.log`.

3. Prepare and execute the cbridge module governance proposal to add the chain:

    1. Prepare the proposal JSON file:

        ```json
        {
          "title": "cBridge config change",
          "description": "Add my BAS",
          "cbr_config": {
            "lp_fee_perc": 100,
            "max_gain_perc": 1200,
            "cbr_contracts": [
              {
                "chain_id": 560,
                "address": "<Bridge-contract-address>"
              }
            ],
            "relay_gas_cost": [
              {
                "chain_id": 560,
                "cost_base": 106987,
                "per_validator": 952,
                "per_sig": 7408
              }
            ],
            "chain_pairs": [
              {
                "chid1": 560,
                "chid2": 56,
                "weight1": 100,
                "fee1to2": 3000,
                "fee2to1": 3000,
                "const_a": 10000
              }
            ],
            "update_gas_cost": true
          },
          "deposit": "0"
        }
        ```

        Important notes:

        - `lp_fee_perc` and `max_gain_perc` are always required. Keep the values as is.
        - `chid1` **MUST BE SMALLER THAN** `chid2`. TODO: Add explanation for params.
        - `deposit` is required. Keep it as "0" for now.

    2. Submit proposal:

        ```sh
        sgnd tx gov submit-proposal cbridge-change <path-to-json-file>
        ```

        Look for a `proposal_id` in the output of the command.

    3. Vote yes on the proposal:

        ```sh
        sgnd tx gov vote <proposal-id> yes
        ```

    4. Wait for the vote period to pass, then check the proposal status:

        ```sh
        sgnd query gov proposal <proposal-id>
        ```

    5. Verify the new cbridge config:

        ```sh
        sgnd query cbridge config
        sgnd query cbridge chaintokens
        ```

4. Reset `Bridge` signers:

    1. Using the `Bridge` contract owner key, call `resetSigners` on the contract with the signer set to the validator
	ETH address and power set to `10000000000000000000000`.

    4. Monitor the `/var/log/sgnd/app.log/` to make sure the `UpdateSigners` event is synced to SGN.

5. Prepare and execute the pegbridge module governance proposal:

    1. Prepare the proposal JSON file:

        ```json
        {
          "title": "pegbridge config change",
          "description": "Add OriginalTokenVault to BSC and PeggedTokenBridge to my BAS",
          "peg_config": {
            "original_token_vaults": [
              {
                "address": "<OriginalTokenVault-contract-address>",
                "chain_id": 56
              }
            ],
            "pegged_token_bridges": [
              {
                "address": "<PeggedTokenBridge-contract-address>",
                "chain_id": 560
              }
            ]
          },
          "deposit": "0"
        }
        ```

    2. Submit proposal:

        ```sh
        sgnd tx gov submit-proposal pegbridge-change <path-to-json-file>
        ```

        Then vote yes on the proposal.

    3. After the voting period has passed, verify the new pegbridge config:

        ```sh
        sgnd query pegbridge config
        ```

### Add pegbridge token

In this example, we assume you want to bridge an original token on BAS to a pegged token on BSC.

1. Deploy pegged token from `sgn-v2-contracts`:

    `MintSwapCanonicalToken` is needed for deploying pegged tokens that follow the Open Canonical Token Standard. Set constructor arguments per the params spreadsheet. Call `updateBridgeSupplyCap` with the `PeggedTokenBridge` address and a cap amount after deployment.

    Edit the env variables in `.env`:

    ```
    MINT_SWAP_CANONICAL_TOKEN_NAME="Your Pegged Token Name"
    MINT_SWAP_CANONICAL_TOKEN_SYMBOL="Your Pegged Token Symbol"
    MINT_SWAP_CANONICAL_TOKEN_DECIMALS=Your Pegged Token Decimals
	  ```

    Then deploy the token:

	  ```sh
    hardhat deploy --tags MintSwapCanonicalToken --network bsc
    ```

2. Add origin token to cbridge module first via a governance proposal:

    1. Prepare proposal JSON:

        ```json
        {
          "title": "cBridge config change",
          "description": "Add FOO on BAS",
          "cbr_config": {
            "lp_fee_perc": 100,
            "max_gain_perc": 1200,
            "assets": [
              {
                "addr": "<FILL-IN>",
                "chain_id": 560,
                "decimal": 18,
                "max_fee_amount": "<FILL-IN>",
                "max_out_amt": "<FILL-IN>",
                "symbol": "FOO",
                "xfer_disabled": false
              }
            ]
          },
          "deposit": "0"
        }
        ```

    2. Submit proposal:

        ```sh
        sgnd tx gov submit-proposal cbridge-change <path-to-json-file>
        ```

        Then vote yes on the proposal.

    3. After the voting period has passed, verify the new cbridge config:

        ```sh
        sgnd query cbridge config
        sgnd query cbridge chaintokens
        ```

3. Add the original-pegged token pair to pegbridge module via a governance proposal:

    1. Prepare the proposal JSON:

        ```json
        {
          "title": "pegbridge config change",
          "description": "Add original and pegged tokens for FOO",
          "peg_config": {
            "orig_pegged_pairs": [
              {
                "burn_fee_pips": 0,
                "max_burn_fee": "<FILL-IN>",
                "max_mint_fee": "<FILL-IN>",
                "mint_fee_pips": 0,
                "orig": {
                  "address": "<original-token-address>",
                  "chain_id": 560,
                  "decimals": 18,
                  "symbol": "FOO"
                },
                "pegged": {
                  "address": "<pegged-token-address>",
                  "chain_id": 56,
                  "decimals": 18,
                  "symbol": "FOO"
                }
              }
            ]
          },
          "deposit": "0"
        }
        ```

    2. Submit proposal:

        ```sh
        sgnd tx gov submit-proposal pegbridge-change <path-to-json-file>
        ```

        Then vote yes on the proposal.

    3. After the voting period has passed, verify the new pegbridge config:

        ```sh
        sgnd query pegbridge config
        ```
