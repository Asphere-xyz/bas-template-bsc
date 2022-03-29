# Run validator with systemd

We recommend using systemd to run your validator. Feel free to experiment with other setups on your own.

1. Prepare the sgnd system service:

    ```sh
    sudo mkdir -p /var/log/sgnd
    sudo touch /var/log/sgnd/tendermint.log
    sudo touch /var/log/sgnd/app.log
    sudo touch /etc/systemd/system/sgnd.service
    ```

    Add the following to `/etc/systemd/system/sgnd.service`:

    ```
    [Unit]
    Description=SGN daemon
    After=network-online.target

    [Service]
    Environment=HOME=/home/ubuntu
    ExecStart=/home/ubuntu/go/bin/sgnd start \
      --home /home/ubuntu/.sgnd/
    StandardOutput=append:/var/log/sgnd/tendermint.log
    StandardError=append:/var/log/sgnd/app.log
    Restart=always
    RestartSec=3
    User=ubuntu
    Group=ubuntu
    LimitNOFILE=4096

    [Install]
    WantedBy=multi-user.target
    ```

    With this setup, `tendermint.log` contains Tendermint logs (i.e. block height, peer info, etc.) while `app.log`
    contains SGN-specific logs.

2. Create `/etc/logrotate.d/sgn` and add the following:

    ```
    /var/log/sgnd/*.log {
        compress
        copytruncate
        daily
        maxsize 30M
        rotate 30
    }
    ```

3. Add an entry to `/etc/crontab` to make logrotate run every 6 hours:

    ```
    30 */6  * * *   root    logrotate -f /etc/logrotate.conf
    ```

4. Enable and start the sgnd service:

    ```sh
    sudo systemctl enable sgnd.service
    sudo systemctl start sgnd.service
    ```

    Now the node should start the state sync. Monitor `tendermint.log` for the progress:

    ```sh
    tail -f /var/log/sgnd/tendermint.log
    ```

    You can tell the node is synced when a new block shows up about every 5 seconds.

## Claim validator status

1. Send 10000 `SGNStakeToken` to your **validator address**, and some ETH for gas to both the **validator and signer addresses**.

2. Initialize the validator by calling [initializeValidator](https://github.com/celer-network/sgn-v2-contracts/blob/448974bb9dd5d2737385a688c4d043ee486e511c/contracts/staking/Staking.sol#L100) on the `Staking` contract and [updateSgnAddr](https://github.com/celer-network/sgn-v2-contracts/blob/448974bb9dd5d2737385a688c4d043ee486e511c/contracts/staking/SGN.sol#L42) on the `SGN` contract.

    Run the init command. Example below sets a `_commissionRate` of 6% and `_minSelfDelegation` (must be equal or greater than the `min self delegations` value obtained at step 2) of 10000 `SGNStakeToken`.

    ```sh
    sgnd ops validator init --commission-rate 0.06 --min-self-delegation 10000 --keystore ~/.sgnd/eth-ks/val.json --passphrase <val-ks-passphrase>
    ```

    Afterwards, verify your validator status:

    ```sh
    sgnd query staking validator <val-eth-address>
    ```

3. Update validator description:

    ```sh
    echo $COSMOS_KEYRING_PASSPHRASE | sgnd tx staking edit-description --website "your-website" --contact "email-address"
    ```

    Note that `COSMOS_KEYRING_PASSPHRASE` here is the passphrase for the keyring used in `sgnd keys add`.

    After a while, verify the updated description:

    ```sh
    sgnd query staking validator <val-eth-address>
    ```

4. Delegate more tokens if your `_minSelfDelegation` is smaller than `min validator tokens` value obtained at step 2.

    ```sh
    sgnd ops delegator delegate --validator <val-eth-address> --amount <amount> --keystore <path-to-keystore-file> --passphrase <ks-passphrase>
    ```

5. Verify validator status. Run following command a few minutes after your validator has enough delegations.

    ```sh
    sgnd query staking validator <val-eth-address>
    ```

    You should see that your validator has the status of `BOND_STATUS_BONDED`.

    You can verify that your validator is in the Tendermint validator set:

    ```sh
    sgnd query tendermint-validator-set
    ```

    You should see an entry with `address` matching the `sgn consensus address` printed as the output of `sgnd ops validator address`.

    You can also verify the delegation:

    ```sh
    sgnd query staking delegation <val-eth-address> <val-eth-address>
    ```

6. (Optional) If something went wrong and your validator is not bonded automatically, you can do so manually through the following command

    ```sh
    sgnd ops validator bond
    ```

    After a while, verify the status:

    ```sh
    sgnd query staking validator <val-eth-address>
    ```

    You should see that your validator has the status of `BOND_STATUS_BONDED`.
