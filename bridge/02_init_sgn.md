# Initialize a Standalone SGN Chain

## Prepare machine and install dependencies

We run on Ubuntu Linux amd64 with Amazon EC2 as an example. Feel free to experiment with other VPS or physical server setups on your own.

1. Start an EC2 machine with the Ubuntu 20.04 LTS image. We recommend using `c6a.2xlarge` if available in the region or `c5.2xlarge` (8 vCPUs, 16GB RAM and 10Gbps network bandwidth) with an EBS volume of at least 500GB. Use the appropriate security groups and a keypair that you have access to.

2. Install go (at least 1.16):

    ```sh
    sudo snap install go --classic
    ```

    Install gcc and make :

    ```sh
    sudo apt update && sudo apt install gcc make
    ```

    (Optional) If you are on Ubuntu and choose to use cleveldb, which is reported to be more performant, install libleveldb-dev:

    ```sh
    sudo apt install libleveldb-dev
    ```

3. Set \$GOBIN and add \$GOBIN to \$PATH. Edit `$HOME/.profile` and add:

    ```sh
    export GOBIN=$HOME/go/bin; export GOPATH=$HOME/go; export PATH=$PATH:$GOBIN
    ```

    to the end, then:

    ```sh
    source $HOME/.profile
    ```

4. (Optional) Install geth:

    ```sh
    sudo add-apt-repository -y ppa:ethereum/ethereum
    sudo apt-get update
    sudo apt-get install ethereum
    ```

## Setup binary, config and accounts

1. From the `/home/ubuntu` directory, download and install the `sgnd` binary:

    ```sh
    curl -L https://github.com/celer-network/sgn-v2-networks/releases/download/v1.6.4/sgnd-v1.6.4-goleveldb-linux-amd64.tar.gz | tar -xz
    # To use with cleveldb on Ubuntu, download https://github.com/celer-network/sgn-v2-networks/releases/download/v1.6.4/sgnd-v1.6.4-cleveldb-ubuntu-linux-amd64.tar.gz
    mv sgnd $GOBIN
    ```

2. From the `bas-template-bsc/bridge` directory, copy the config files:

    ```sh
    mkdir -p $HOME/.sgnd/config
    cp sgn-configs $HOME/.sgnd/config
    ```

3. Initialize the new validator node:

    `node-name` is a name you specify for the node. `network-name` can be something like `sgn-bas-1`.

    ```sh
	# Remove existing genesis.json first
    rm $HOME/.sgnd/config/genesis.json
	# Initialize default genesis.json and config.toml
    sgnd init <node-name> --chain-id <network-name> --home $HOME/.sgnd
    # Overwrite genesis.json and config.toml with the ones from sgn-v2-networks
    cp genesis.json config.toml $HOME/.sgnd/config
    # Create an empty Tendermint snapshots directory
    mkdir -p $HOME/.sgnd/data/snapshots
    ```

    Backup the generated Tendermint key files `$HOME/.sgnd/config/node_key.json` and `$HOME/.sgnd/config/priv_validator_key.json` securely. Make sure the keys are **never** committed to any repo.

4. Fill out the fields in the Tendermint config file `$HOME/.sgnd/config/config.toml` with the correct values:

    | Field | Description |
    | ----- | ----------- |
    | moniker | The `node-name` you decided |
    | external_address| `<public-ip:26656>`, where `public-ip` is the public IP of the machine hosting the node |
    | db_backend | `goleveldb` or `cleveldb` depending on the binary used |

5. Add a Cosmos SDK / Tendermint validator account:

    ```sh
    sgnd keys add <node-name>
    ```

    Input a passphrase for the keyring. Backup the passphrase along with the displayed mnemonic phrase securely. Make sure they are
	**never** committed to any repo.

    To view the account created, run:

    ```sh
    sgnd keys list
    ```

    Make a note of the **sgn-prefixed account address**.

6. Prepare an Ethereum key as the **validator key**:

    Assuming `geth` is installed, the keystore JSON file can be generated via:

    ```sh
    geth account new
    ```

    Backup the passphrase securely. Save the JSON file as `$HOME/.sgnd/eth-ks/val.json`:

    ```sh
    mkdir $HOME/.sgnd/eth-ks
    cp <path-to-keystore-json> $HOME/.sgnd/eth-ks/val.json
    ```

7. Prepare another Ethereum key as the **signer key**, which will be used for signing cross-chain transactions. Follow the same steps as
    preparing the validator key. Save the JSON file as `$HOME/.sgnd/eth-ks/signer.json`.

8. Fill out the fields in the SGN-specific config file `$HOME/.sgnd/config/sgn.toml` with the correct values:

    | Field | Description |
    | ----- | ----------- |
    | eth.gateway | Your BAS gateway URL |
    | eth.signer_keystore | The path to the signer Ethereum keystore file in step 7 |
    | eth.signer_passphrase | The passphrase of the signer keystore |
    | eth.validator_address | The **Ethereum address** of the validator key prepared in step 6 |
    | sgnd.passphrase | The **Cosmos keyring passphrase** you typed in step 5 |
    | sgnd.validator_account | The **sgn-prefixed validator Cosmos SDK account** added in step 5 |

9. Fill in the missing gateway URLs in `$HOME/.sgnd/config/cbridge.toml` with the corresponding JSON-RPC URLs for the chains.

10. Fill in the fields with the required info in `$HOME/.sgnd/config/genesis.json`. Note that we use 560 as the BAS chain ID but feel free
    to change it to your own. TODO: Add more explanations of the fields.
